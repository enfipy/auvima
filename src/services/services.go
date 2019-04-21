package services

import (
	"context"
	"errors"
	"net"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"

	videoController "github.com/enfipy/auvima/src/services/video/controller"
	videoDeliveryCron "github.com/enfipy/auvima/src/services/video/delivery/cron"
	videoDeliveryHttp "github.com/enfipy/auvima/src/services/video/delivery/http"
	videoUsecase "github.com/enfipy/auvima/src/services/video/usecase"

	"github.com/enfipy/locker"
	"github.com/robfig/cron"
)

func InitServices(cnfg *config.Config) (start func(list net.Listener), close func()) {
	if cnfg.Settings == nil {
		helpers.PanicOnError(errors.New("Valid settings must be provided"))
	}

	cronInstance := cron.New()
	locker := locker.Initialize()
	echo := helpers.InitHttp()
	pc := helpers.InitPostgres()
	coubClient := helpers.InitCoubClient(
		cnfg.Settings.Coub.Urls.BaseUrl,
		cnfg.Settings.Coub.Secrets.AccessToken,
	)
	instaClient := helpers.InitInstagramClient(
		cnfg.Settings.Instagram.Creds.Username,
		cnfg.Settings.Instagram.Creds.Password,
		cnfg.Settings.Instagram.CredsPath,
	)
	youtubeClient := helpers.InitYoutubeClient(
		cnfg.Settings.Youtube.Creds.AccessToken,
		cnfg.Settings.Youtube.Creds.TokenType,
		cnfg.Settings.Youtube.Creds.RefreshToken,
		cnfg.Settings.Youtube.Creds.Expiry,
	)

	videoUcs := videoUsecase.NewUsecase(cnfg, pc, locker)
	videoCnr := videoController.NewController(cnfg, videoUcs, coubClient, instaClient, youtubeClient)
	videoDeliveryHttp.NewHttp(echo, cnfg, videoCnr)
	videoDeliveryCron.NewCron(cronInstance, cnfg, videoCnr)

	start = func(lis net.Listener) {
		echo.Listener = lis
		cronInstance.Start()
		echo.Start("")
	}
	close = func() {
		cronInstance.Stop()
		echo.Shutdown(context.Background())
	}
	return
}
