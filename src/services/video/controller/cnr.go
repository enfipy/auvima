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

func (cnr *videoController) SaveCoub(permalink string) {
	coub := cnr.GetCoub(permalink)

	mp4Path, mp4Link := cnr.GetMediaInfo(coub.Permalink, "mp4", &coub.FileVersions.HTML5.Video)
	rmVideo := DownloadCoub(mp4Path, mp4Link)

	mp3Path, mp3Link := cnr.GetMediaInfo(coub.Permalink, "mp3", &coub.FileVersions.HTML5.Audio)
	rmAudio := DownloadAudio(mp3Path, mp3Link)

	go func() {
		defer rmVideo()
		defer rmAudio()
		cnr.SaveFinishedVideo(mp4Path, mp3Path, coub)
	}()
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
	if coub.Duration < 10 {
		loopTimes = 2
	}

	duration := fmt.Sprintf("%f", coub.Duration*float64(loopTimes))
	out := fmt.Sprintf("%s/%s.mp4", cnr.config.Settings.Storage.Finished, coub.Permalink)
	loop := fmt.Sprintf("loop=%d:size=32767:start=0", loopTimes)

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
