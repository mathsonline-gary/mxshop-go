package initialize

import (
	"fmt"

	"mxshop-go/user_api/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Config() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("user_api")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Debugf("Configs: %v", global.ServerConfig)

	// Watcher
	viper.OnConfigChange(func(evt fsnotify.Event) {
		fmt.Println("config changed: " + evt.Name)
		_ = viper.ReadInConfig()
		_ = viper.Unmarshal(global.ServerConfig)
	})
	viper.WatchConfig()
}