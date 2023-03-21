package tools

import (
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/********************************************************************
created:    2020-08-14
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var exitCallbacks ExitCallbacks
var signalChan = make(chan os.Signal)

type ExitCallbacks struct {
	sync.Mutex
	d []func(os.Signal)
}

func (my *ExitCallbacks) Register(callback func(sig os.Signal)) {
	if callback != nil {
		my.Lock()
		my.d = append(my.d, callback)
		my.Unlock()
	}
}

func (my *ExitCallbacks) Execute(sig os.Signal) {
	my.Lock()
	defer my.Unlock()

	for i := len(my.d) - 1; i >= 0; i-- {
		var callback = my.d[i]
		callback(sig)
	}
}

func InitBasics() {
	// 设置时区:东八，这个经常会报data race问题
	//time.Local, _ = time.LoadLocation("Asia/Chongqing")

	// 监听 Ctrl+C, kill
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	// 给rand随机化一个种子
	rand.Seed(time.Now().UnixNano())

	go func() {
		select {
		case sig := <-signalChan:
			// 不知道为什么，这个时灵时不灵的，愁人
			//fmt.Printf("exiting application by signal=%q, len(exitCallbacks)=%d\n", sig, len(exitCallbacks))
			exitCallbacks.Execute(sig)
			os.Exit(0)
		}
	}()

	// 设置默认ENV
	const key = "ENV"
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, "DEV")
	}
}

// 注册 exit callback 回调
func RegisterExitCallback(callback func(sig os.Signal)) {
	exitCallbacks.Register(callback)
}
