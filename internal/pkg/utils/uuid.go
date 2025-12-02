package utils

import "github.com/google/uuid"

func UUID() uuid.UUID {
	res, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return res
}
