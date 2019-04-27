package usecase

import (
	"time"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/models"
	"github.com/enfipy/auvima/src/services/video"

	"github.com/enfipy/locker"
	"github.com/jackc/pgx/pgtype"
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

func (ucs *videoUsecase) SaveVideo(uniqueId string, duration int64, origin models.VideoOrigin) *models.Video {
	video := &models.Video{
		UniqueId:  uniqueId,
		Duration:  duration,
		Status:    models.VideoStatus_Normal,
		Origin:    origin,
		CreatedAt: time.Now().Unix(),
	}

	ucs.pc.Exec(`
		INSERT INTO videos(unique_id, duration, status, origin, created_at)
		VALUES($1, $2, $3, $4, $5)
		`, video.UniqueId, video.Duration, video.Status, video.Origin, video.CreatedAt,
	)

	return video
}

func (ucs *videoUsecase) GetUnusedVideos(limit int32) []models.Video {
	rows := ucs.pc.QueryMany(`
		SELECT id, unique_id, duration, used_in, status, origin, created_at
		FROM videos
		WHERE used_in IS NULL
		ORDER BY RANDOM()
		LIMIT $1
	`, limit)
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		video := models.Video{}
		var usedIn pgtype.Int4

		rows.Scan(&video.Id, &video.UniqueId, &video.Duration, &usedIn, &video.Status, &video.Origin, &video.CreatedAt)
		if video.Id == 0 {
			continue
		}

		if usedIn.Int == 0 {
			id := uint32(usedIn.Int)
			video.UsedIn = &id
		}

		videos = append(videos, video)
	}

	return videos
}

func (ucs *videoUsecase) GetVideosByUsedIn(id uint32) []models.Video {
	rows := ucs.pc.QueryMany(`
		SELECT id, unique_id, duration, used_in, status, origin, created_at
		FROM videos
		WHERE used_in = $1
	`, id)
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		video := models.Video{}
		var usedIn pgtype.Int4

		rows.Scan(&video.Id, &video.UniqueId, &video.Duration, &usedIn, &video.Status, &video.Origin, &video.CreatedAt)
		if video.Id == 0 {
			continue
		}

		if usedIn.Int == 0 {
			id := uint32(usedIn.Int)
			video.UsedIn = &id
		}

		videos = append(videos, video)
	}

	return videos
}

func (ucs *videoUsecase) UseVideo(uniqueId string, productionVideoId uint32) {
	ucs.pc.Exec(`
		UPDATE videos
		SET used_in = $2
		WHERE unique_id = $1
		`, uniqueId, productionVideoId,
	)
}

func (ucs *videoUsecase) UseProduction(id uint32) {
	ucs.pc.Exec(`
		UPDATE prods
		SET used = TRUE
		WHERE id = $1
		`, id,
	)
}

func (ucs *videoUsecase) GetVideo(uniqueId string) *models.Video {
	getResult := ucs.pc.Query(`
		SELECT id, unique_id, used_in, duration, status, origin, created_at
		FROM videos
		WHERE unique_id = $1
	`, uniqueId)

	video := models.Video{}
	var usedIn pgtype.Int4
	getResult(&video.Id, &video.UniqueId, &usedIn, &video.Duration, &video.Status, &video.Origin, &video.CreatedAt)
	if video.Id == 0 {
		return nil
	}

	if usedIn.Int == 0 {
		id := uint32(usedIn.Int)
		video.UsedIn = &id
	}

	return &video
}

func (ucs *videoUsecase) SaveProd(duration int64) *models.Production {
	prod := &models.Production{
		Duration:  duration,
		Used:      false,
		CreatedAt: time.Now().Unix(),
	}

	ucs.pc.Exec(`
		INSERT INTO prods(duration, used, created_at)
		VALUES($1, $2, $3)
		`, prod.Duration, prod.Used, prod.CreatedAt,
	)

	return prod
}

func (ucs *videoUsecase) GetProductionVideoCount() int64 {
	getResult := ucs.pc.Query(`
		SELECT COUNT(id)
		FROM prods
	`)

	var count pgtype.Int8
	getResult(&count)

	return count.Int
}

func (ucs *videoUsecase) GetNextProductionVideo() *models.Production {
	getResult := ucs.pc.Query(`
		SELECT id, duration, used, created_at
		FROM prods
		WHERE id = (
			SELECT max(id)
			FROM prods
			WHERE used = false
		)
	`)

	prod := models.Production{}
	getResult(&prod.Id, &prod.Duration, &prod.Used, &prod.CreatedAt)

	if prod.Id == 0 {
		return nil
	}
	return &prod
}
