package utils

import (
	"encoding/hex"

	"github.com/google/uuid"
)

func UUID() uuid.UUID {
	res, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return res
}

func UniqueID() string {
	id := UUID()
	return hex.EncodeToString(id[:])
}
