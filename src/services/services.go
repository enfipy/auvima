package services

import (
	"errors"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"

	videoController "github.com/enfipy/auvima/src/services/video/controller"
	videoDeliveryHttp "github.com/enfipy/auvima/src/services/video/delivery/http"
	videoUsecase "github.com/enfipy/auvima/src/services/video/usecase"

	"github.com/enfipy/locker"
	"github.com/labstack/echo"
)

func InitServices(cnfg *config.Config) (*echo.Echo, func()) {
	if cnfg.Settings == nil {
		helpers.PanicOnError(errors.New("Valid settings must be provided"))
	}

	locker := locker.Initialize()
	echo := helpers.InitHttp()
	pc := helpers.InitPostgres()
	coubClient := helpers.InitCoubClient(cnfg.Settings.Coub.Urls.BaseUrl, cnfg.Settings.Coub.Secrets.AccessToken)

	videoUcs := videoUsecase.NewUsecase(cnfg, pc, locker)
	videoCnr := videoController.NewController(cnfg, videoUcs, coubClient)
	videoDeliveryHttp.NewDelivery(echo, videoCnr)

	return echo, func() {}
}
