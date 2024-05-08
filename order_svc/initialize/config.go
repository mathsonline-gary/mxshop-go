package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"github.com/zycgary/mxshop-go/order_svc/config"
	"github.com/zycgary/mxshop-go/order_svc/global"
)

func initConfig() {
	fmt.Println("configurations initializing...")

	viper.SetConfigName("grpc")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/order")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Get config driver
	if err := viper.UnmarshalKey("dc", &global.Config.DC); err != nil {
		panic(err)
	}

	switch global.Config.DC.Driver {
	case "nacos":
		var nacosConfig config.Nacos
		if err := viper.UnmarshalKey("nacos.server", &nacosConfig.Server); err != nil {
			panic(err)
		}
		if err := viper.UnmarshalKey("nacos.client", &nacosConfig.Client); err != nil {
			panic(err)
		}

		// Fetch configs from Nacos
		nacosClientConfig := global.Config.Nacos.Client
		nacosServerConfig := global.Config.Nacos.Server
		clientConfig := constant.ClientConfig{
			NamespaceId:         nacosClientConfig.Namespace,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "tmp/nacos/log",
			CacheDir:            "tmp/nacos/cache",
			LogLevel:            "debug",
		}
		serverConfigs := []constant.ServerConfig{
			{
				IpAddr:      nacosServerConfig.Host,
				Port:        nacosServerConfig.Port,
				ContextPath: "/nacos",
				Scheme:      "http",
			},
		}
		configClient, err := clients.CreateConfigClient(map[string]interface{}{
			"serverConfigs": serverConfigs,
			"clientConfig":  clientConfig,
		})
		if err != nil {
			panic(err)
		}

		content, err := configClient.GetConfig(vo.ConfigParam{
			DataId: nacosClientConfig.DataID,
			Group:  nacosClientConfig.Group,
		})
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal([]byte(content), &global.Config); err != nil {
			panic(err)
		}

		// Avoid overwriting of associate configs
		global.Config.DC.Driver = "nacos"
		global.Config.Nacos = nacosConfig

	case "local":
		if err := viper.Unmarshal(&global.Config); err != nil {
			panic(err)
		}
	}

	fmt.Println("configurations initialized!")
}
