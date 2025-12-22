package app

import (
	"context"

	"github.com/Miklakapi/go-file-share/internal/config"
	"github.com/Miklakapi/go-file-share/internal/files"
)

type DependencyBag struct {
	Config     config.Config
	FileHub    *files.FileHub
	AppContext context.Context
}

func NewDependencyBag(config config.Config, fileHub *files.FileHub, appContext context.Context) *DependencyBag {
	return &DependencyBag{
		Config:     config,
		FileHub:    fileHub,
		AppContext: appContext,
	}
}
