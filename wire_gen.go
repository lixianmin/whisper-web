// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/lixianmin/whisper-web/app/builds"
	"github.com/lixianmin/whisper-web/app/web"
	"github.com/lixianmin/whisper-web/app/web/console"
)

// Injectors from wire.ignore.go:

func buildApplication() *builds.Application {
	viper := builds.NewConfig()
	iLogger := builds.InitLogger(viper)
	engine := web.NewEngine(viper, iLogger)
	ginServeMux := web.NewGinServeMux(engine)
	app := console.NewRoadApp(viper)
	server := console.NewConsole(viper, ginServeMux, app)
	application := builds.NewApplication(viper, server)
	return application
}