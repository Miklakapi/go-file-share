package app

import (
	"context"

	"github.com/Miklakapi/go-file-share/internal/config"
	fileShareService "github.com/Miklakapi/go-file-share/internal/file-share/application"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type DependencyBag struct {
	AppContext context.Context
	Config     config.Config

	RoomRepo  ports.RoomRepository
	FileStore ports.FileStore
	Hasher    ports.PasswordHasher

	FileShareService *fileShareService.Service
}

func NewDependencyBag(
	appContext context.Context,
	config config.Config,
	roomRepo ports.RoomRepository,
	fileStore ports.FileStore,
	hasher ports.PasswordHasher,
) *DependencyBag {
	fileShareService := fileShareService.NewService(roomRepo, fileStore, hasher)

	return &DependencyBag{
		AppContext: appContext,
		Config:     config,

		RoomRepo:  roomRepo,
		FileStore: fileStore,
		Hasher:    hasher,

		FileShareService: fileShareService,
	}
}
