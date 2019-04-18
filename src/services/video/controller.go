package video

import "github.com/enfipy/auvima/src/models"

type Controller interface {
	GetCoub(permalink string) *models.Coub
	SaveCoub(coub *models.Coub) *models.Video
	GetCoubs(tag, order string, page, perPage int) []models.Coub
}
