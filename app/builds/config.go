package builds

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

/********************************************************************
created:    2020-05-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func NewConfig() *viper.Viper {
	//环境变量
	var env, ok = os.LookupEnv("ENV")
	if !ok {
		env = "DEV"
	}

	var vip = viper.New()
	vip.SetEnvPrefix("TOUR")
	vip.AddConfigPath("res/config/")

	var configName = "app_" + strings.ToLower(env)
	vip.SetConfigName(configName)

	err := vip.ReadInConfig() // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return vip
}
