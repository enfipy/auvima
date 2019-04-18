package models

import "github.com/google/uuid"

type Video struct {
	Id        uuid.UUID
	UniqueId  string
	Used      bool
	CreatedAt int64
}
