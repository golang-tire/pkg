package utiles

import "github.com/google/uuid"

// IsValidUUID return true if given uuid is valid
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
