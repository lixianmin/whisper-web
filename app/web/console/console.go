package console

import (
	"fmt"
	"github.com/lixianmin/gonsole"
	"github.com/lixianmin/gonsole/road"
	"github.com/lixianmin/gonsole/road/epoll"
	"github.com/lixianmin/got/timex"
	"github.com/lixianmin/whisper-web/app/ifs"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

/********************************************************************
created:    2020-07-04
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type (
	StringTable = map[string]interface{}
)

func NewConsole(vip *viper.Viper, mux gonsole.IServeMux, roadApp *road.App) *gonsole.Server {
	//获取web的端口和密码
	var webPort = vip.GetInt(ifs.ConfigServerWebPort)
	var userPasswords = vip.GetStringMapString("server.userPasswords")

	var processName = filepath.Base(os.Args[0])
	var title = fmt.Sprintf("%s %s", os.Getenv("ENV"), processName)

	var console = gonsole.NewServer(mux,
		gonsole.WithAutoLoginTime(7*timex.Day),
		gonsole.WithPort(webPort),
		gonsole.WithPageTitle(title),
		gonsole.WithUserPasswords(userPasswords),
		gonsole.WithEnablePProf(true),
		gonsole.WithDeadlockIgnores([]string{
			"internal/poll.runtime_pollWait(",
			"time.Sleep(",
		}),
		// 大的日志下载不下来，跟这个参数没关系，使用127的ip在本机下载很多，但是使用180.76.162.195的
		// 外网ip下载很快就会被中断，怀疑是BCC的网络设置有问题
		//WriteBufferSize: 1024 * 1024,
	)

	return console
}

func NewRoadApp(vip *viper.Viper) *road.App {
	var tcpPort = vip.GetString(ifs.ConfigServerTcpPort)
	var tcpAddress = fmt.Sprintf(":%s", tcpPort)

	var acceptor = epoll.NewTcpAcceptor(tcpAddress)
	var app = road.NewApp(acceptor,
		// 限流，原来是2，太小了，改大一些
		road.WithSessionRateLimitBySecond(1000),
	)

	return app
}
