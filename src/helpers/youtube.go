package helpers

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func InitYoutubeClient(accessToken, tokenType, refreshToken string, expiry uint64) *http.Client {
	ctx := context.Background()
	tok := &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		Expiry:       time.Now().Add(time.Duration(expiry) * time.Second),
		RefreshToken: refreshToken,
	}
	cnfg := oauth2.Config{}
	client := cnfg.Client(ctx, tok)
	return client
}
