package logger

import (
	"errors"
	"github.com/lixianmin/logo"
	"github.com/lixianmin/logo/ding"
	"github.com/lixianmin/logo/lark"
	"github.com/lixianmin/whisper-web/core/tools"
	"os"
	"strings"
)

/********************************************************************
created:    2020-05-26
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var theLogger = logo.GetLogger().(*logo.Logger)

type InitArgs struct {
	FilterLevel string
	Lark        *lark.Lark
}

func Init(args InitArgs) {
	// 开启异步写标记，提高日志输出性能
	theLogger.AddFlag(logo.LogAsyncWrite)

	// 调整theLogger的filterLevel
	var level = getFilterLevel(args.FilterLevel)
	theLogger.SetFilterLevel(level)

	// 文件日志
	const flag = logo.FlagDate | logo.FlagTime | logo.FlagShortFile | logo.FlagLevel
	var rollingFile = logo.NewRollingFileHook(logo.RollingFileHookArgs{Flag: flag, FilterLevel: level})
	theLogger.AddHook(rollingFile)

	if !tools.IsDev() && args.Lark != nil {
		// 飞书报警
		var hook = ding.NewHook(args.Lark, ding.WithFilterLevel(logo.LevelWarn))
		theLogger.AddHook(hook)
	}

	// 这里故意不判nil，因为这个参数不参为nil
	tools.RegisterExitCallback(func(sig os.Signal) {
		// 程序退出时关闭logger以及所有实现了Closer接口的appenders
		theLogger.Info("[Init()] logger is exiting...")
		_ = theLogger.Close()
	})

	//// 设置depth=2，输出一条日志后，为了方便外面调用，再设置depth=3
	//theLogger.SetFuncCallDepth(2)
	//theLogger.Info("[Init()] log init done")
	//theLogger.SetFuncCallDepth(3)
}

func getFilterLevel(filterLevel string) int {
	var level = strings.ToLower(filterLevel)
	switch level {
	case "debug":
		return logo.LevelDebug
	case "warn", "warning":
		return logo.LevelWarn
	case "error":
		return logo.LevelError
	default:
		return logo.LevelInfo
	}
}

func GetLogger() *logo.Logger {
	return theLogger
}

func Dot(err interface{}) error {
	if err != nil {
		switch err := err.(type) {
		case string:
			var v = errors.New(err)
			theLogger.Error(err)
			return v
		case error:
			theLogger.Error("err=%q", err)
			return err
		}
	}

	return nil
}
