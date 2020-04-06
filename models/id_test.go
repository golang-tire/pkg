package models

import (
	"github.com/google/uuid"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	// test valid uuid
	validUUID := GenerateUUID()
	_, err := uuid.Parse(validUUID)
	if err != nil{
		t.Error("GenerateUUID() generate invalid uuid")
	}
}
