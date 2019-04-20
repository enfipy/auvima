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

func (ucs *videoUsecase) SaveVideo(uniqueId string, origin models.VideoOrigin) *models.Video {
	video := &models.Video{
		Id:        uuid.New(),
		UniqueId:  uniqueId,
		Used:      false,
		Status:    models.VideoStatus_Normal,
		Origin:    origin,
		CreatedAt: time.Now().Unix(),
	}

	ucs.pc.Exec(`
		INSERT INTO videos
		VALUES($1, $2, $3, $4, $5, $6)
		`, video.Id, video.UniqueId, video.Used, video.Status, video.Origin, video.CreatedAt,
	)

	return video
}

func (ucs *videoUsecase) GetUnusedVideos(limit int32) []models.Video {
	rows := ucs.pc.QueryMany(`
		SELECT id, unique_id, used, status, origin, created_at
		FROM videos
		WHERE used = FALSE
		ORDER BY RANDOM()
		LIMIT $1
	`, limit)
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		video := models.Video{}

		rows.Scan(&video.Id, &video.UniqueId, &video.Used, &video.Status, &video.Origin, &video.CreatedAt)
		if !helpers.ValidateUUID(video.Id) {
			continue
		}

		videos = append(videos, video)
	}

	return videos
}

func (ucs *videoUsecase) UseVideo(uniqueId string) {
	ucs.pc.Exec(`
		UPDATE videos
		SET used = TRUE
		WHERE unique_id = $1
		`, uniqueId,
	)
}
