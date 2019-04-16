package helpers

import (
	"github.com/badoux/checkmail"
	"github.com/google/uuid"
)

func isAlphanumeric(char rune) bool {
	if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') {
		return false
	}
	return true
}

func ValidateUUID(id interface{}) bool {
	if id == nil {
		return false
	}
	parsedId, ok := id.(uuid.UUID)
	if !ok {
		tryId, ok := id.([16]byte)
		if !ok {
			return false
		}
		parsedId = uuid.UUID(tryId)
	}
	if parsedId.String() == "00000000-0000-0000-0000-000000000000" {
		return false
	}
	return true
}

func ParseUUID(id string) (*uuid.UUID, bool) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, false
	}
	if parsedId.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, false
	}
	return &parsedId, true
}

func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return false
	}
	return true
}
