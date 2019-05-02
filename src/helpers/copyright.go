package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	AuddResponse struct {
		Status string  `json:"status"`
		Result *Result `json:"result"`
		Error  *Error  `json:"error"`
	}

	Error struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}

	Result struct {
		Artist      string `json:"artist"`
		Title       string `json:"title"`
		Album       string `json:"album"`
		ReleaseDate string `json:"release_date"`
		Label       string `json:"label"`
	}
)

func CheckCopyright() *AuddResponse {
	source := "https://youtu.be/glEiPXAYE-U"
	link := fmt.Sprintf("http://api.audd.io/?return=timecode,itunes,deezer,lyrics&itunes_country=us&url=%s", source)
	resp, err := http.Get(link)
	PanicOnError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	PanicOnError(err)

	var res AuddResponse
	err = json.Unmarshal(body, &res)
	PanicOnError(err)

	return &res
}
