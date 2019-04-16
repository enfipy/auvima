package delivery

import (
	"errors"

	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/services/video"

	echoHTTP "github.com/labstack/echo"
)

type videoServer struct {
	videoController video.Controller
}

func NewDelivery(echo *echoHTTP.Echo, videoController video.Controller) {
	server := &videoServer{
		videoController: videoController,
	}

	echo.GET("/api/v1/video/coub", helpers.Handle(server.GetCoub))
}

func (server *videoServer) GetCoub(ctx echoHTTP.Context) interface{} {
	permalink := ctx.QueryParam("permalink")
	if permalink == "" {
		panic(errors.New("Permalink must be provided"))
	}

	coub := server.videoController.GetCoub(permalink)
	return coub
}
