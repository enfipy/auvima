package http

import (
	"errors"
	"log"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/models"
	"github.com/enfipy/auvima/src/services/video"

	echoHTTP "github.com/labstack/echo"
)

type videoServer struct {
	config          *config.Config
	videoController video.Controller
}

func NewHttp(echo *echoHTTP.Echo, config *config.Config, videoController video.Controller) {
	server := &videoServer{
		config:          config,
		videoController: videoController,
	}

	echo.GET("/api/v1/video/coub", helpers.Handle(server.SaveCoub))
	echo.GET("/api/v1/video/coubs", helpers.Handle(server.GetCoubs))
	echo.GET("/api/v1/video/instagram", helpers.Handle(server.GetInstagramVideos))
	echo.GET("/api/v1/video/gen", helpers.Handle(server.GenerateVideo))
	echo.GET("/api/v1/video/upload", helpers.Handle(server.UploadVideo))
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

func (server *videoServer) GetInstagramVideos(_ echoHTTP.Context) interface{} {
	var videos []models.Video
	for _, account := range server.config.Settings.Instagram.MaterialAccounts {
		log.Print("Loading videos from " + account)
		createdVideos := server.videoController.GetInstagramVideos(account, 3)
		videos = append(videos, createdVideos...)
	}
	return videos
}

func (server *videoServer) GenerateVideo(_ echoHTTP.Context) interface{} {
	server.videoController.GenerateProductionVideo()
	return nil
}

func (server *videoServer) UploadVideo(_ echoHTTP.Context) interface{} {
	id := server.videoController.UploadVideo()
	return id
}
