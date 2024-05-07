package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"github.com/zycgary/mxshop-go/stock_svc/global"
)

func initConfig() {
	fmt.Println("configurations initializing...")

	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("stock_svc")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Get Nacos configs
	if err := viper.UnmarshalKey("nacos.server", &global.NacosConfig.NacosServerConfig); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("nacos.client", &global.NacosConfig.NacosClientConfig); err != nil {
		panic(err)
	}

	// Fetch configs from Nacos
	nacosClientConfig := global.NacosConfig.NacosClientConfig
	nacosServerConfig := global.NacosConfig.NacosServerConfig
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

	fmt.Println("configurations initialized!")
}
