package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

func InitCoubClient(uri, accessToken string) *CoubClient {
	baseURL, err := url.Parse(uri)
	PanicOnError(err)

	var cnfg *oauth2.Config
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	client := &CoubClient{
		Client:  cnfg.Client(oauth2.NoContext, token),
		BaseURL: baseURL,
	}

	return client
}

type CoubClient struct {
	Client  *http.Client
	BaseURL *url.URL
}

func (client *CoubClient) NewRequest(method, urlString string, body interface{}) *http.Request {
	rel, err := url.Parse(urlString)
	PanicOnError(err)

	uri := client.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		PanicOnError(err)
	}

	req, err := http.NewRequest(method, uri.String(), buf)
	PanicOnError(err)

	return req
}

func (client *CoubClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
