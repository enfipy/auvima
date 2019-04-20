package controller

import (
	"fmt"
	"os/exec"

	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/models"
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

func (cnr *videoController) ScaleVideo(uniqueId, url string) {
	filter := "[0:0]split[main][back];" +
		"[back]scale=1920:1080[scale];" +
		"[scale]drawbox=x=0:y=0:w=1920:h=1080:color=black:t=1000[draw];" +
		"[main]scale='if(gt(a,16/9),1920,-1)':'if(gt(a,16/9),-1,1080)'[proc];" +
		"[draw][proc]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2[fhd]; [fhd]setsar=1/1[sarfix]"

	path := fmt.Sprintf("%s/%s.mp4", cnr.config.Settings.Storage.Temporary, uniqueId)
	rmFile := DownloadFile(path, url)
	defer rmFile()

	out := GetPath(cnr.config.Settings.Storage.Finished, uniqueId)
	commandArgs := []string{
		"-i", path,
		"-filter_complex", filter,
		"-map", "[sarfix]", "-map", "0",
		"-y", out,
	}

	cmd := exec.Command("ffmpeg", commandArgs...)
	err := cmd.Run()
	helpers.PanicOnError(err)
}

func (cnr *videoController) ScaleAndLoopVideo(mp4Path, mp3Path, uniqueId string, dur float64, loopTimes int) {

	duration := fmt.Sprintf("%f", dur)
	out := GetPath(cnr.config.Settings.Storage.Finished, uniqueId)
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

func (cnr *videoController) ConcatVideo(videos []models.Video, op, frame25, end string) {
	commandArgs := []string{"-i", op}
	for _, video := range videos {
		path := GetPath(cnr.config.Settings.Storage.Finished, video.UniqueId)
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

	// Todo: Make special quality for youtube
	out := GetPath(cnr.config.Settings.Storage.Production, name)
	filter := fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", inputCount)

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
}
