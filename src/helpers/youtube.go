package helpers

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	youtube "google.golang.org/api/youtube/v3"
)

func InitYoutubeClient(accessToken, tokenType, refreshToken string, expiry uint64) *youtube.Service {
	ctx := context.Background()
	tok := &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		Expiry:       time.Now().Add(time.Duration(expiry) * time.Second),
		RefreshToken: refreshToken,
	}
	cnfg := oauth2.Config{}
	client := cnfg.Client(ctx, tok)

	service, err := youtube.New(client)
	PanicOnError(err)

	return service
}
