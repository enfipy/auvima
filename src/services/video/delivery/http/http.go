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

	echo.GET("/api/v1/video/coub", helpers.Handle(server.SaveCoub))
	echo.GET("/api/v1/video/coubs", helpers.Handle(server.GetCoubs))
	echo.GET("/api/v1/video/gen", helpers.Handle(server.GenerateVideo))
}

func (server *videoServer) SaveCoub(ctx echoHTTP.Context) interface{} {
	permalink := ctx.QueryParam("permalink")
	if permalink == "" {
		panic(errors.New("Permalink must be provided"))
	}

	coub := server.videoController.GetCoub(permalink)
	server.videoController.SaveCoub(coub)
	return coub
}

func (server *videoServer) GetCoubs(ctx echoHTTP.Context) interface{} {
	tag := ctx.QueryParam("tag")
	if tag == "" {
		panic(errors.New("Tag must be provided"))
	}

	coubs := server.videoController.GetCoubs(tag, "newest_popular", 1, 10)
	for _, coub := range coubs {
		server.videoController.SaveCoub(&coub)
	}
	return coubs
}

func (server *videoServer) GenerateVideo(_ echoHTTP.Context) interface{} {
	server.videoController.GenerateProductionVideo()
	return nil
}
