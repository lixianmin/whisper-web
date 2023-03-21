package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixianmin/logo"
	"github.com/lixianmin/whisper-web/app/ifs"
	"github.com/lixianmin/whisper-web/core/logger"
	"github.com/spf13/viper"
)

/********************************************************************
created:    2020-06-12
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const (
	keyData = "data" // gin中，context中作为返回值的data
)

func NewEngine(vip *viper.Viper, log logo.ILogger) *gin.Engine {
	var engine = newGin()

	// 初始化gonsole
	var webPort = vip.GetInt(ifs.ConfigServerWebPort)
	registerUrls(engine)

	var addr = fmt.Sprintf(":%d", webPort)
	var certFile = "res/ssl/localhost.crt"
	var keyFile = "res/ssl/localhost.key"
	go engine.RunTLS(addr, certFile, keyFile)

	log.Info("starting gin, webPort=%d", webPort)
	return engine
}

func registerUrls(engine *gin.Engine) {
	//var rootPath = "api/v1"
	//var rootGroup = engine.Group(rootPath, postInterceptor)
	//{
	//	var tools = rootGroup.Group("/tools")
	//	{
	//		registerGet(tools, "/ip", controller.GetToolsIP)
	//	}
	//}
}

func registerGet(group *gin.RouterGroup, relativePath string, handler func(ctx *gin.Context) (interface{}, error)) {
	group.GET(relativePath, func(context *gin.Context) {
		var data, err = handler(context)
		if err != nil {
			_ = context.Error(err)
			logger.GetLogger().Info(err.Error())
		} else {
			context.Set(keyData, data)
		}
	})
}

func newGin() *gin.Engine {
	engine := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout

	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			if param.Latency > time.Minute {
				// Truncate in a golang < 1.8 safe way
				param.Latency = param.Latency - param.Latency%time.Second
			}

			return fmt.Sprintf("[%s %12v] %s %#v\n%s",
				param.ClientIP,
				param.Latency,
				param.Method,

				param.Path,
				param.ErrorMessage,
			)
		},
		Output: logger.GetLogger(),
	}))
	engine.Use(gin.Recovery())

	return engine
}

func postInterceptor(context *gin.Context) {
	context.Next()

	var result struct {
		ErrNo  int         `json:"errno"`
		ErrMsg string      `json:"errmsg"`
		Data   interface{} `json:"data"`
	}

	// context.Errors有时候为nil，有时候不是，但len()=0
	if len(context.Errors) == 0 {
		result.Data, _ = context.Get(keyData)
	} else {
		result.ErrNo = http.StatusInternalServerError
		result.ErrMsg = fmt.Sprintf("%+v", context.Errors)
		result.Data = nil
	}

	context.JSON(http.StatusOK, result)
}
