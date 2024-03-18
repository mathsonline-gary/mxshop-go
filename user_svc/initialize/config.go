package initialize

import (
	"fmt"

	"mxshop-go/user_svc/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("user_svc")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	// Watcher
	viper.OnConfigChange(func(evt fsnotify.Event) {
		fmt.Println("config changed: " + evt.Name)
		_ = viper.ReadInConfig()
		_ = viper.Unmarshal(global.ServerConfig)
	})
	viper.WatchConfig()
}
