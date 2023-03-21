package main

import (
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/whisper-web/core/tools"
)

/********************************************************************
created:    2023-03-21
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func main() {
	defer loom.DumpIfPanic()
	tools.InitBasics()
	var application = buildApplication()
	application.Run()
}
