package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadCoub(path, link string) func() {
	out, err := os.Create(path)
	PanicOnError(err)
	defer out.Close()

	resp, err := http.Get(link)
	PanicOnError(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	PanicOnError(err)

	// Decode weird coub encoding
	_, err = out.WriteAt([]byte{0}, 0)
	PanicOnError(err)
	_, err = out.WriteAt([]byte{0}, 1)
	PanicOnError(err)

	return func() {
		os.Remove(path)
	}
}

func DownloadFile(path, link string) func() {
	out, err := os.Create(path)
	PanicOnError(err)
	defer out.Close()

	resp, err := http.Get(link)
	PanicOnError(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	PanicOnError(err)

	return func() {
		os.Remove(path)
	}
}

func GetPath(storagePath, fileName string) string {
	return fmt.Sprintf("%s/%s.mp4", storagePath, fileName)
}
