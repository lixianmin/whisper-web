package builds

import (
	"os"
	"runtime"

	"github.com/google/wire"
	jsonitor "github.com/json-iterator/go"
	"github.com/lixianmin/gonsole"
	"github.com/lixianmin/got/convert"
	"github.com/lixianmin/logo"
	"github.com/lixianmin/whisper-web/app/web"
	"github.com/lixianmin/whisper-web/app/web/console"
	"github.com/lixianmin/whisper-web/core/tools"
	"github.com/spf13/viper"
)

/********************************************************************
created:    2020-05-26
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 给wire系统提供provider
// 1. 不建议放到wire.go文件中，因为wire.go不参与编译，因此无法借助golang自动导入头文件
// 2. 每次修改完成后，在项目根目录下执行wire命令，重新生成wire_gen.go文件

var BuildProviderSet = wire.NewSet(InitLogger,
	NewConfig,
	web.NewEngine,
	wire.Bind(new(gonsole.IServeMux), new(*web.GinServeMux)),
	web.NewGinServeMux,
	console.NewConsole,
	console.NewRoadApp,
	NewApplication)

type Application struct {
	vip *viper.Viper
}

func NewApplication(vip *viper.Viper, console *gonsole.Server) *Application {
	convert.InitJson(jsonitor.Marshal, jsonitor.Unmarshal)
	tools.RegisterExitCallback(func(sig os.Signal) {
		logo.Info("[NewApplication()] exiting application by signal=%q", sig)
	})

	logo.JsonW("goVersion", runtime.Version(), "gitBranchName", gonsole.GitBranchName, "gitCommitId", gonsole.GitCommitId, "gitCommitMessage",
		gonsole.GitCommitMessage, "gitCommitTime", gonsole.GitCommitTime, "appBuildTime", gonsole.AppBuildTime, "console", console.ConsoleUrl())

	return &Application{
		vip: vip,
	}
}

func (my *Application) Run() {
	logo.Debug("Application Run...")
	select {}
}
