package controller

import (
	"fmt"
	"os"

	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/models"

	goinsta "github.com/ahmdrz/goinsta/v2"
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

	return durations[0]
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

	return durations[0]
}

func (cnr *videoController) ConcatVideo(videos []models.Video, op, end, frame25 string) (string, int64) {
	commandArgs := []string{"-i", op}
	for _, video := range videos {
		path := helpers.GetPath(cnr.config.Settings.Storage.Finished, video.UniqueId)
		commandArgs = append(commandArgs, "-i", path, "-i", frame25)
	}
	commandArgs = commandArgs[:len(commandArgs)-2]
	commandArgs = append(commandArgs, "-i", end)

	inputCount := 0
	for _, str := range commandArgs {
		if str == "-i" {
			inputCount++
		}
	}

	// Todo: Get name based on unique ids
	name, err := helpers.GenRandString(5)
	helpers.PanicOnError(err)

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

	durations := helpers.GetDurations(output)
	if len(durations) <= 0 {
		helpers.PanicOnError(err)
	}
	helpers.PanicOnError(err)

	return name, durations[0]
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
