package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/models"
	"github.com/enfipy/auvima/src/services/video"

	goinsta "github.com/ahmdrz/goinsta/v2"
)

type videoController struct {
	config       *config.Config
	videoUsecase video.Usecase

	coubClient  *helpers.CoubClient
	instaClient *goinsta.Instagram
}

func NewController(
	cnfg *config.Config, videoUsecase video.Usecase, coubClient *helpers.CoubClient, instaClient *goinsta.Instagram,
) video.Controller {
	return &videoController{
		config:       cnfg,
		videoUsecase: videoUsecase,
		coubClient:   coubClient,
		instaClient:  instaClient,
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
	rmAudio := DownloadFile(mp3Path, mp3Link)
	defer rmAudio()

	cnr.SaveFinishedVideo(mp4Path, mp3Path, coub)
	video := cnr.videoUsecase.SaveVideo(coub.Permalink, models.VideoOrigin_Coub)

	return video
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

	cnr.ScaleAndLoopVideo(mp4Path, mp3Path, coub.Permalink, dur, loopTimes)
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

func (cnr *videoController) GenerateProductionVideo() {
	videos := cnr.videoUsecase.GetUnusedVideos(10)
	if len(videos) <= 0 {
		panic(errors.New("Can not generate video. No prepared videos"))
	}

	op := GetPath(cnr.config.Settings.Storage.Static, "op")
	end := GetPath(cnr.config.Settings.Storage.Static, "end")
	frame25 := GetPath(cnr.config.Settings.Storage.Static, "25frame")

	cnr.ConcatVideo(videos, op, end, frame25)

	// Todo: Update video. Set used = true
}

func (cnr *videoController) GetInstagramVideos(username string, limit int) []models.Video {
	user, err := cnr.instaClient.Profiles.ByName(username)
	helpers.PanicOnError(err)

	sh := time.Duration(cnr.config.Settings.Instagram.SuitabilityHours)
	timestamp := time.Now().Add(-sh * time.Hour).Unix()
	from := strconv.FormatInt(timestamp, 10)

	videos := GetVideosFromInstagramUser(user, from, limit)

	var savedVideos []models.Video
	for uniqueId, url := range videos {
		cnr.ScaleVideo(uniqueId, url)
		video := cnr.videoUsecase.SaveVideo(uniqueId, models.VideoOrigin_Instagram)
		savedVideos = append(savedVideos, *video)
	}

	return savedVideos
}
