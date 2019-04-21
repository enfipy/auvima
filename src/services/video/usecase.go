package video

import "github.com/enfipy/auvima/src/models"

type Usecase interface {
	SaveVideo(uniqueId string, duration int64, origin models.VideoOrigin) *models.Video
	SaveProd(duration int64) *models.Production

	UseVideo(uniqueId string, productionVideoId uint32)
	UseProduction(id uint32)

	GetUnusedVideos(limit int32) []models.Video
	GetVideosByUsedIn(id uint32) []models.Video
	GetVideo(uniqueId string) *models.Video
	GetProductionVideoCount() int64
	GetNextProductionVideo() *models.Production
}
