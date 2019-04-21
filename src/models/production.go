package models

import "github.com/google/uuid"

type Production struct {
	Id        uuid.UUID
	UniqueId  string
	Duration  int64
	CreatedAt int64
}
