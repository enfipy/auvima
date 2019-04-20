package models

import "github.com/google/uuid"

type Video struct {
	Id        uuid.UUID
	UniqueId  string
	Used      bool
	Status    VideoStatus
	Origin    VideoOrigin
	CreatedAt int64
}

type VideoStatus uint8

const (
	VideoStatus_Normal VideoStatus = iota
	VideoStatus_Urgent
)

type VideoOrigin uint8

const (
	VideoOrigin_Instagram VideoOrigin = iota
	VideoOrigin_Coub
)
