package initialize

import (
	"fmt"

	"mxshop-go/user_api/global/config"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Config() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("user_api")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Debugf("Configs: %v", config.ServerConfig)

	// Watcher
	viper.OnConfigChange(func(evt fsnotify.Event) {
		fmt.Println("config changed: " + evt.Name)
		_ = viper.ReadInConfig()
		_ = viper.Unmarshal(config.ServerConfig)
	})
	viper.WatchConfig()
}
