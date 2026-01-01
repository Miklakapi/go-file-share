package app

import (
	"github.com/Miklakapi/go-file-share/internal/config"
	fileShare "github.com/Miklakapi/go-file-share/internal/file-share/application"
	fileShareDomain "github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type DependencyBag struct {
	Config config.Config

	RoomRepo     ports.RoomRepository
	FileStore    ports.FileStore
	Hasher       ports.PasswordHasher
	TokenService ports.TokenService

	FileShareService *fileShare.Service
}

func NewDependencyBag(
	config config.Config,
	roomRepo ports.RoomRepository,
	fileStore ports.FileStore,
	hasher ports.PasswordHasher,
	tokenService ports.TokenService,
) *DependencyBag {
	fileShareSettings := fileShareDomain.NewPolicy(
		config.DefaultRoomTTL,
		config.TokenTTL,
		config.MaxFiles,
		config.MaxRoomBytes,
		config.MaxRoomLifespan,
		config.MaxTokenLifespan,
		config.UploadDir,
	)

	fileShareService := fileShare.NewService(roomRepo, fileStore, hasher, tokenService, fileShareSettings)

	return &DependencyBag{
		Config: config,

		RoomRepo:     roomRepo,
		FileStore:    fileStore,
		Hasher:       hasher,
		TokenService: tokenService,

		FileShareService: fileShareService,
	}
}
