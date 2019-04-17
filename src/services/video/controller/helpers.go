package controller

import (
	"io"
	"net/http"
	"os"

	"github.com/enfipy/auvima/src/helpers"
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

func DownloadAudio(path, link string) func() {
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
