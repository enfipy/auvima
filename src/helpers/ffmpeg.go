package helpers

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

const (
	Milisecond = 1
	Second     = 1000 * Milisecond
	Minute     = 60 * Second
	Hour       = 60 * Minute

	DurationString = "Duration: "
	StartString    = ", start"
)

func ExecFFMPEG(commandArgs []string) ([]string, error) {
	cmd := exec.Command("ffmpeg", commandArgs...)

	var bytes bytes.Buffer
	cmd.Stderr = &bytes

	err := cmd.Run()
	output := bytes.String()
	lines := strings.Split(output, "\n")

	if err != nil {
		errorMessageIndex := len(lines) - 2
		if errorMessageIndex < 0 {
			return lines, err
		}
		errorMessage := lines[errorMessageIndex]
		return lines, errors.New(errorMessage)
	}

	return lines, nil
}

func GetDurations(lines []string) []int64 {
	var durations []int64
	for _, line := range lines {
		if strings.Contains(line, DurationString) {
			time := getBetween(line)

			duration := strings.Split(time, ":")
			if len(duration) < 3 {
				continue
			}

			hour, _ := strconv.ParseInt(duration[0], 10, 64)
			min, _ := strconv.ParseInt(duration[1], 10, 64)
			secFloat, _ := strconv.ParseFloat(duration[2], 64)

			sec := int64(secFloat * Second)

			result := hour*Hour + min*Minute + sec
			durations = append(durations, result)
		}
	}
	return durations
}

func getBetween(line string) string {
	posFirst := strings.Index(line, DurationString)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(line, StartString)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(DurationString)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return line[posFirstAdjusted:posLast]
}
