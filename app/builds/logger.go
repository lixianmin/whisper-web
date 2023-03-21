package builds

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lixianmin/got/osx"
	"github.com/lixianmin/logo"
	"github.com/lixianmin/logo/lark"
	"github.com/lixianmin/whisper-web/app/ifs"
	"github.com/lixianmin/whisper-web/core/logger"
	"github.com/spf13/viper"
)

/********************************************************************
created:    2021-01-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 初始化日志
func InitLogger(vip *viper.Viper) logo.ILogger {
	var webPort = vip.GetInt(ifs.ConfigServerWebPort)
	var title = fetchTitlePrefix(webPort)

	var token = strings.TrimSpace(vip.GetString("logger.larkToken"))
	var talk *lark.Lark
	if token != "" {
		talk = lark.NewLark(title, token)
	}

	// 初始化日志
	logger.Init(logger.InitArgs{
		FilterLevel: vip.GetString("logger.filterLevel"),
		Lark:        talk,
	})

	logo.Info("logger inited, title=%q\n\n", title)
	return logger.GetLogger()
}

func fetchTitlePrefix(port int) string {
	var processName = filepath.Base(os.Args[0])
	var localIP = osx.GetLocalIp()
	var prefix = fmt.Sprintf("%s %s %s:%d", os.Getenv("ENV"), processName, localIP, port)
	return prefix
}
