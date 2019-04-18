package usecase

import (
	"time"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/models"
	"github.com/enfipy/auvima/src/services/video"

	"github.com/enfipy/locker"
	"github.com/google/uuid"
)

type videoUsecase struct {
	config *config.Config
	pc     *helpers.PostgresConnection
	locker *locker.Locker
}

func NewUsecase(cnfg *config.Config, pc *helpers.PostgresConnection, locker *locker.Locker) video.Usecase {
	return &videoUsecase{
		config: cnfg,
		pc:     pc,
		locker: locker,
	}
}

func (ucs *videoUsecase) SaveCoub(permalink string) *models.Video {
	video := &models.Video{
		Id:        uuid.New(),
		UniqueId:  permalink,
		Used:      false,
		CreatedAt: time.Now().Unix(),
	}

	ucs.pc.Exec(`
		INSERT INTO videos
		VALUES($1, $2, $3, $4)
		`, video.Id, video.UniqueId, video.Used, video.CreatedAt,
	)

	return video
}
