package models

type Video struct {
	Id uint32

	UniqueId string
	UsedIn   *uint32
	Duration int64

	Status VideoStatus
	Origin VideoOrigin

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
