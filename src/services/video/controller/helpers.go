package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/enfipy/auvima/src/helpers"

	goinsta "github.com/ahmdrz/goinsta/v2"
)

func DownloadCoub(path, link string) func() {
	out, err := os.Create(path)
	helpers.PanicOnError(err)
	defer out.Close()

	resp, err := http.Get(link)
	helpers.PanicOnError(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	helpers.PanicOnError(err)

	// Decode weird coub encoding
	_, err = out.WriteAt([]byte{0}, 0)
	helpers.PanicOnError(err)
	_, err = out.WriteAt([]byte{0}, 1)
	helpers.PanicOnError(err)

	return func() {
		os.Remove(path)
	}
}

func DownloadFile(path, link string) func() {
	out, err := os.Create(path)
	helpers.PanicOnError(err)
	defer out.Close()

	resp, err := http.Get(link)
	helpers.PanicOnError(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	helpers.PanicOnError(err)

	return func() {
		os.Remove(path)
	}
}

func GetPath(storagePath, fileName string) string {
	return fmt.Sprintf("%s/%s.mp4", storagePath, fileName)
}

func GetVideosFromInstagramUser(user *goinsta.User, from string, limit int) map[string]string {
	media := user.Feed(from)
	videos := map[string]string{}
	for media.Next() {
		for _, item := range media.Items {
			if len(item.Videos) != 0 {
				for _, itemVideo := range item.Videos {
					if len(videos) >= limit {
						return videos
					}

					videos[item.Code] = itemVideo.URL
				}
			}
		}
	}
	return videos
}
