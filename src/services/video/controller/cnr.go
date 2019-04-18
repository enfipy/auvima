package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/models"
	"github.com/enfipy/auvima/src/services/video"
)

type videoController struct {
	config       *config.Config
	videoUsecase video.Usecase
	coubClient   *helpers.CoubClient
}

func NewController(cnfg *config.Config, videoUsecase video.Usecase, coubClient *helpers.CoubClient) video.Controller {
	return &videoController{
		config:       cnfg,
		videoUsecase: videoUsecase,
		coubClient:   coubClient,
	}
}

func (cnr *videoController) GetCoub(permalink string) *models.Coub {
	link := "api/v2/coubs/" + permalink

	req := cnr.coubClient.NewRequest("GET", link, nil)
	resp, _ := cnr.coubClient.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	helpers.PanicOnError(err)

	var res models.Coub
	err = json.Unmarshal(body, &res)
	helpers.PanicOnError(err)

	return &res
}

func (cnr *videoController) SaveCoub(coub *models.Coub) *models.Video {
	mp4Path, mp4Link := cnr.GetMediaInfo(coub.Permalink, "mp4", &coub.FileVersions.HTML5.Video)
	rmVideo := DownloadCoub(mp4Path, mp4Link)
	defer rmVideo()

	mp3Path, mp3Link := cnr.GetMediaInfo(coub.Permalink, "mp3", &coub.FileVersions.HTML5.Audio)
	rmAudio := DownloadAudio(mp3Path, mp3Link)
	defer rmAudio()

	cnr.SaveFinishedVideo(mp4Path, mp3Path, coub)
	video := cnr.videoUsecase.SaveVideo(coub.Permalink)

	return video
}

func (cnr *videoController) GetMediaInfo(permalink, format string, media *models.Media) (path, link string) {
	if media.High != nil {
		link = media.High.URL
	} else {
		link = media.Med.URL
	}
	path = fmt.Sprintf("%s/%s.%s", cnr.config.Settings.Storage.Temporary, permalink, format)
	return
}

func (cnr *videoController) SaveFinishedVideo(mp4Path, mp3Path string, coub *models.Coub) {
	loopTimes := 1
	dur := coub.Duration
	if coub.Duration <= 5 {
		loopTimes = 3
		dur = dur * float64(loopTimes)
	} else if coub.Duration > 20 {
		dur = 10
	}

	duration := fmt.Sprintf("%f", dur)
	out := GetPath(cnr.config.Settings.Storage.Finished, coub.Permalink)
	filter := fmt.Sprintf("[0:0]split[main][back];"+
		"[back]scale=1920:1080[scale];"+
		"[scale]drawbox=x=0:y=0:w=1920:h=1080:color=black:t=1000[draw];"+
		"[main]scale='if(gt(a,16/9),1920,-1)':'if(gt(a,16/9),-1,1080)'[proc];"+
		"[draw][proc]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2[fhd]; [fhd]setsar=1/1[sarfix];"+
		"[sarfix]loop=%d:size=9999:start=0",
		loopTimes,
	)

	commandArgs := []string{
		"-i", mp4Path,
		"-i", mp3Path,
		"-filter_complex", filter,
		"-map", "0", "-map", "1",
		"-t", duration,
		"-y", out,
	}

	cmd := exec.Command("ffmpeg", commandArgs...)

	err := cmd.Run()
	helpers.PanicOnError(err)
}

func (cnr *videoController) GetCoubs(_, order string, page, perPage int) []models.Coub {
	var res struct {
		Page       int           `json:"page"`
		PerPage    int           `json:"per_page"`
		TotalPages int           `json:"total_pages"`
		Coubs      []models.Coub `json:"coubs"`
	}

	// link := fmt.Sprintf("http://coub.com/api/v2/timeline/tag/%s?page=%d&per_page=%d&order_by=%s", tag, page, perPage, order)
	link := fmt.Sprintf("http://coub.com/api/v2/timeline/hot?page=%d&per_page=%d&order_by=%s", page, perPage, order)

	req := cnr.coubClient.NewRequest("GET", link, nil)
	resp, _ := cnr.coubClient.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	helpers.PanicOnError(err)

	err = json.Unmarshal(body, &res)
	helpers.PanicOnError(err)

	return res.Coubs
}

func (cnr *videoController) GenerateProductionVideo() {
	videos := cnr.videoUsecase.GetUnusedVideos(10)
	if len(videos) <= 0 {
		panic(errors.New("Can not generate video. No prepared videos"))
	}

	op := GetPath(cnr.config.Settings.Storage.Static, "op")
	end := GetPath(cnr.config.Settings.Storage.Static, "end")
	frame25 := GetPath(cnr.config.Settings.Storage.Static, "25frame")

	commandArgs := []string{"-i", op}
	for _, video := range videos {
		path := GetPath(cnr.config.Settings.Storage.Finished, video.UniqueId)
		commandArgs = append(commandArgs, "-i", path, "-i", frame25)
	}
	commandArgs = commandArgs[:len(commandArgs)-2]
	commandArgs = append(commandArgs, "-i", end)

	// Todo: Fix this ugly place
	count := 0
	for _, str := range commandArgs {
		if str == "-i" {
			count++
		}
	}

	// Todo: Get name based on unique ids
	name, err := helpers.GenRandString(5)
	helpers.PanicOnError(err)

	// Todo: Make special quality for youtube
	out := GetPath(cnr.config.Settings.Storage.Production, name)
	filter := fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", count)

	commandArgs = append(
		commandArgs,
		"-map", "[v]",
		"-map", "[a]",
		"-filter_complex", filter,
		"-y", out,
	)
	cmd := exec.Command("ffmpeg", commandArgs...)

	err = cmd.Run()
	helpers.PanicOnError(err)

	// Todo: Update video. Set used = true
}
