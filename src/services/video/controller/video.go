package controller

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/models"

	goinsta "github.com/ahmdrz/goinsta/v2"
	youtube "google.golang.org/api/youtube/v3"
)

func (cnr *videoController) GetMediaInfo(permalink, format string, media *models.Media) (path, link string) {
	if media.High != nil {
		link = media.High.URL
	} else {
		link = media.Med.URL
	}
	path = fmt.Sprintf("%s/%s.%s", cnr.config.Settings.Storage.Temporary, permalink, format)
	return
}

func (cnr *videoController) ScaleVideo(uniqueId, url string) int64 {
	filter := "[0:0]split[main][back];" +
		"[back]scale=1920:1080[scale];" +
		"[scale]drawbox=x=0:y=0:w=1920:h=1080:color=black:t=1000[draw];" +
		"[main]scale='if(gt(a,16/9),1920,-1)':'if(gt(a,16/9),-1,1080)'[proc];" +
		"[draw][proc]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2[fhd];" +
		"[fhd]setsar=1/1[sarfix]"

	path := helpers.GetPath(cnr.config.Settings.Storage.Temporary, uniqueId)
	rmFile := helpers.DownloadFile(path, url)
	defer rmFile()

	out := helpers.GetPath(cnr.config.Settings.Storage.Finished, uniqueId)
	commandArgs := []string{
		"-i", path,
		"-filter_complex", filter,
		"-map", "[sarfix]", "-map", "0:1",
		"-vcodec", "libx264",
		"-filter:a", "loudnorm",
		"-y", out,
	}
	output, err := helpers.ExecFFMPEG(commandArgs)

	durations := helpers.GetDurations(output)
	if len(durations) <= 0 {
		helpers.PanicOnError(err)
	}
	helpers.PanicOnError(err)

	return durations[len(durations)-1]
}

func (cnr *videoController) ScaleAndLoopVideo(mp4Path, mp3Path, uniqueId string, dur float64, loopTimes int) int64 {
	duration := fmt.Sprintf("%f", dur)
	out := helpers.GetPath(cnr.config.Settings.Storage.Finished, uniqueId)
	filter := fmt.Sprintf("[0:0]split[main][back];"+
		"[back]scale=1920:1080[scale];"+
		"[scale]drawbox=x=0:y=0:w=1920:h=1080:color=black:t=1000[draw];"+
		"[main]scale='if(gt(a,16/9),1920,-1)':'if(gt(a,16/9),-1,1080)'[proc];"+
		"[draw][proc]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2[fhd];"+
		"[fhd]setsar=1/1[sarfix];"+
		"[sarfix]loop=%d:size=9999:start=0",
		loopTimes,
	)

	commandArgs := []string{
		"-i", mp4Path,
		"-i", mp3Path,
		"-filter_complex", filter,
		"-map", "0", "-map", "1",
		"-vcodec", "libx264",
		"-filter:a", "loudnorm",
		"-t", duration,
		"-y", out,
	}
	output, err := helpers.ExecFFMPEG(commandArgs)

	durations := helpers.GetDurations(output)
	if len(durations) <= 0 {
		helpers.PanicOnError(err)
	}
	helpers.PanicOnError(err)

	return durations[len(durations)-1]
}

func (cnr *videoController) ConcatVideo(videos []models.Video, op, end, frame25 string) int64 {
	var duration int64 = cnr.config.Settings.Video.OpLength + cnr.config.Settings.Video.FrameLength
	commandArgs := []string{"-i", op}

	for _, video := range videos {
		path := helpers.GetPath(cnr.config.Settings.Storage.Finished, video.UniqueId)
		commandArgs = append(commandArgs, "-i", path, "-i", frame25)
		duration += video.Duration
	}
	commandArgs = commandArgs[:len(commandArgs)-2]
	commandArgs = append(commandArgs, "-i", end)

	inputCount := 0
	for _, str := range commandArgs {
		if str == "-i" {
			inputCount++
		}
	}

	productionVideoCount := cnr.videoUsecase.GetProductionVideoCount()
	if productionVideoCount <= 0 {
		productionVideoCount = 1
	} else {
		productionVideoCount++
	}
	name := strconv.FormatInt(productionVideoCount, 10)

	out := helpers.GetPath(cnr.config.Settings.Storage.Production, name)
	filter := fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", inputCount)

	commandArgs = append(
		commandArgs,
		"-filter_complex", filter,
		"-map", "[v]",
		"-map", "[a]",
		"-vcodec", "libx264",
		"-y", out,
	)
	output, err := helpers.ExecFFMPEG(commandArgs)
	helpers.PanicOnError(err)

	durations := helpers.GetDurations(output)
	if len(durations) <= 0 {
		helpers.PanicOnError(errors.New("Fatal: No durations in output"))
	}

	return duration
}

func (cnr *videoController) GetVideosFromInstagramUser(user *goinsta.User, from string, limit int) map[string]string {
	media := user.Feed(from)
	videos := map[string]string{}
	for media.Next() {
		for _, item := range media.Items {
			if len(item.Videos) != 0 {
				for _, itemVideo := range item.Videos {
					if len(videos) >= limit {
						return videos
					}

					uniqueId := item.Code

					existedVideo := cnr.videoUsecase.GetVideo(uniqueId)
					if existedVideo != nil {
						continue
					}

					videos[uniqueId] = itemVideo.URL
				}
			}
		}
	}
	return videos
}

func (cnr *videoController) RemoveVideo(uniqueId string) {
	path := helpers.GetPath(cnr.config.Settings.Storage.Finished, uniqueId)
	err := os.Remove(path)
	helpers.PanicOnError(err)
}

func (cnr *videoController) PrepareYoutubeVideo(videoNumber uint32, videos []models.Video) (*youtube.Video, string) {
	videoCnfg := cnr.config.Settings.Video

	title := fmt.Sprintf(videoCnfg.Title, videoNumber)
	filename := helpers.GetPath(
		cnr.config.Settings.Storage.Production,
		strconv.FormatInt(int64(videoNumber), 10),
	)

	var links string
	var dur int64 = cnr.config.Settings.Video.OpLength + cnr.config.Settings.Video.FrameLength
	for _, video := range videos {
		timeM := dur / helpers.Minute
		timeS := (dur % helpers.Minute) / helpers.Second

		if video.Origin == models.VideoOrigin_Coub {
			links += fmt.Sprintf("%d:%d - https://coub.com/view/%s\n", timeM, timeS, video.UniqueId)
		}
		if video.Origin == models.VideoOrigin_Instagram {
			links += fmt.Sprintf("%d:%d - https://www.instagram.com/p/%s\n", timeM, timeS, video.UniqueId)
		}

		dur += video.Duration + cnr.config.Settings.Video.FrameLength
	}
	description := fmt.Sprintf(videoCnfg.Description, links)

	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  videoCnfg.CategoryId,
			Tags:        strings.Split(videoCnfg.Tags, ","),
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: videoCnfg.Privacy,
		},
	}

	return video, filename
}
