package video

import "github.com/enfipy/auvima/src/models"

type Usecase interface {
	SaveVideo(uniqueId string, duration int64, origin models.VideoOrigin) *models.Video
	SaveProd(uniqueId string, duration int64) *models.Production

	UseVideo(uniqueId string)

	GetUnusedVideos(limit int32) []models.Video
	GetVideo(uniqueId string) *models.Video
}
