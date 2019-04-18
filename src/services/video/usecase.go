package video

import "github.com/enfipy/auvima/src/models"

type Usecase interface {
	SaveCoub(permalink string) *models.Video
}
