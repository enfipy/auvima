package video

import "github.com/enfipy/auvima/src/models"

type Usecase interface {
	SaveVideo(uniqueId string, origin models.VideoOrigin) *models.Video
	GetUnusedVideos(limit int32) []models.Video
	UseVideo(uniqueId string)
}
