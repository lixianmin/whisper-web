package web

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lixianmin/got/osx"
)

/********************************************************************
created:    2020-06-12
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

//GinServeMux是自己定义的结构体，包含engine结构体
type GinServeMux struct {
	engine *gin.Engine
}

func NewGinServeMux(engine *gin.Engine) *GinServeMux {
	var logRoot = "logs"
	var mux = &GinServeMux{engine: engine}
	mux.handleLogs(logRoot)

	return mux
}

// 因为gin的特点，它处理日志时，其pattern必须是"/logs/:name"这样的格式
func (mux *GinServeMux) handleLogs(logRoot string) {
	var pattern = fmt.Sprintf("/%s/:name", logRoot)
	mux.engine.GET(pattern, func(context *gin.Context) {
		var logFilePath = context.Request.URL.Path
		if len(logFilePath) < 1 {
			return
		}

		logFilePath = logFilePath[1:]
		if osx.IsPathExist(logFilePath) {
			var bytes, err = ioutil.ReadFile(logFilePath)
			if err == nil {
				_, err = context.Writer.Write(bytes)
			} else {
				var text = fmt.Sprintf("err=%q", err)
				_, err = context.Writer.Write([]byte(text))
			}

			if err != nil {
				fmt.Println(err)
			}
		}
	})
}

func (mux *GinServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.engine.GET(pattern, func(context *gin.Context) {
		handler(context.Writer, context.Request)
	})
}
