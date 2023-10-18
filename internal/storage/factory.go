package storage

import (
	"context"
)

type ClientManager interface {
	Download(destFolder string, filepath string)
}

type Storage struct {
	context context.Context
}

func GetStorage(kind string) ClientManager {
	if kind == "do_spaces" {
		return newDOStorage()
	}

	// default PublicStorage
	return newPublicStorage()
}
