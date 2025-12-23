package app

import (
	"context"

	"github.com/Miklakapi/go-file-share/internal/config"
	fileShare "github.com/Miklakapi/go-file-share/internal/file-share/application"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type DependencyBag struct {
	AppContext context.Context
	Config     config.Config

	RoomRepo     ports.RoomRepository
	FileStore    ports.FileStore
	Hasher       ports.PasswordHasher
	TokenService ports.TokenService

	FileShareService *fileShare.Service
}

func NewDependencyBag(
	appContext context.Context,
	config config.Config,
	roomRepo ports.RoomRepository,
	fileStore ports.FileStore,
	hasher ports.PasswordHasher,
	tokenService ports.TokenService,
) *DependencyBag {
	fileShareSettings := fileShare.NewSettings(
		config.DefaultRoomTTL,
		config.TokenTTL,
		config.MaxFiles,
		config.MaxRoomBytes,
		config.MaxRoomLifespan,
		config.MaxTokenLifespan,
		config.CleanupInterval,
	)

	fileShareService := fileShare.NewService(roomRepo, fileStore, hasher, tokenService, fileShareSettings)

	return &DependencyBag{
		AppContext: appContext,
		Config:     config,

		RoomRepo:     roomRepo,
		FileStore:    fileStore,
		Hasher:       hasher,
		TokenService: tokenService,

		FileShareService: fileShareService,
	}
}
