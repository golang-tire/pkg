package models

import "github.com/google/uuid"

// GenerateUUID generates a unique UUID that can be used as an identifier for an entity.
func GenerateUUID() string {
	return uuid.New().String()
}
