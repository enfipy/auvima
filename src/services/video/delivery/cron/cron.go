package cron

import (
	"log"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/services/video"

	"github.com/robfig/cron"
)

type videoServer struct {
	config          *config.Config
	cronInstance    *cron.Cron
	videoController video.Controller
}

func NewCron(cronInstance *cron.Cron, config *config.Config, videoController video.Controller) {
	server := &videoServer{
		config:          config,
		cronInstance:    cronInstance,
		videoController: videoController,
	}

	timingsCnfg := server.config.Settings.Video.Timings

	server.cronInstance.AddFunc(timingsCnfg.FetchMaterial, server.FetchMaterial)
	server.cronInstance.AddFunc(timingsCnfg.GenerateProductionVideo, server.GenerateProductionVideo)
	server.cronInstance.AddFunc(timingsCnfg.UploadVideo, server.UploadVideo)
}

func (server *videoServer) FetchMaterial() {
	defer helpers.RecoverWithLog()

	instCnfg := server.config.Settings.Instagram
	for _, account := range instCnfg.MaterialAccounts {
		log.Print("Loading videos from " + account)
		limit := int(instCnfg.MaterialCountToFetch)
		createdVideos := server.videoController.GetInstagramVideos(account, limit)
		log.Printf("Loaded %d videos", len(createdVideos))
	}
}

func (server *videoServer) GenerateProductionVideo() {
	defer helpers.RecoverWithLog()

	log.Print("Generating production video")
	prod := server.videoController.GenerateProductionVideo()
	log.Printf("Generated production video with %s id", prod.UniqueId)
}

func (server *videoServer) UploadVideo() {
	defer helpers.RecoverWithLog()

	log.Print("Uploading production video")
	server.videoController.UploadVideo()
	log.Print("Video uploaded successfully")
}
