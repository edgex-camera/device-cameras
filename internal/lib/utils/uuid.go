package utils

import (
	uuid "github.com/satori/go.uuid"
)

func GenUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}
