package controller

import (
	"encoding/json"
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
	video := cnr.videoUsecase.SaveCoub(coub.Permalink)

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
	if coub.Duration < 5 {
		loopTimes = 3
	} else if coub.Duration < 10 {
		loopTimes = 2
		dur = dur * float64(loopTimes)
	} else if coub.Duration > 20 {
		dur = 10
	}

	duration := fmt.Sprintf("%f", dur)
	out := fmt.Sprintf("%s/%s.mp4", cnr.config.Settings.Storage.Finished, coub.Permalink)
	loop := fmt.Sprintf("loop=%d:size=9999:start=0", loopTimes)

	cmd := exec.Command(
		"ffmpeg",
		"-i", mp4Path,
		"-i", mp3Path,
		"-t", duration,
		"-filter_complex", loop,
		"-y", out,
	)

	err := cmd.Run()
	helpers.PanicOnError(err)
}

func (cnr *videoController) GetCoubs(tag, order string, page, perPage int) []models.Coub {
	var res struct {
		Page       int           `json:"page"`
		PerPage    int           `json:"per_page"`
		TotalPages int           `json:"total_pages"`
		Coubs      []models.Coub `json:"coubs"`
	}

	link := fmt.Sprintf("http://coub.com/api/v2/timeline/tag/%s?page=%d&per_page=%d&order_by=%s", tag, page, perPage, order)

	req := cnr.coubClient.NewRequest("GET", link, nil)
	resp, _ := cnr.coubClient.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	helpers.PanicOnError(err)

	err = json.Unmarshal(body, &res)
	helpers.PanicOnError(err)

	return res.Coubs
}
