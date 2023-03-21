package tools

import "os"

/********************************************************************
created:    2020-06-18
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func IsDevOrFat() bool {
	return IsDev() || IsFat()
}

func IsDev() bool {
	var env = os.Getenv("ENV")
	return env == "DEV"
}

func IsFat() bool {
	var env = os.Getenv("ENV")
	return env == "FAT"
}

func IsPro() bool {
	var env = os.Getenv("ENV")
	return env == "PRO"
}
