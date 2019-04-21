package video

import "github.com/enfipy/auvima/src/models"

type Usecase interface {
	SaveVideo(uniqueId string, duration int64, origin models.VideoOrigin) *models.Video
	GetUnusedVideos(limit int32) []models.Video
	GetVideo(uniqueId string) *models.Video
	UseVideo(uniqueId string)
}
