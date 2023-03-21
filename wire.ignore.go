//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/lixianmin/whisper-web/app/builds"
)

func buildApplication() *builds.Application {
	wire.Build(builds.BuildProviderSet)
	return &builds.Application{}
}
