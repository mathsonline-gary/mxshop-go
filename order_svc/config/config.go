package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/spf13/viper"
)

type App struct {
	Name  string `mapstructure:"name" json:"name"`
	Host  string `mapstructure:"host" json:"host"`
	Port  int32  `mapstructure:"port" json:"port"`
	Env   string `mapstructure:"env" json:"env"`
	Debug bool   `mapstructure:"debug" json:"debug"`
}

type Log struct {
	Channel string `mapstructure:"channel" json:"channel"`
}

type DB struct {
	Connection      string `mapstructure:"connection" json:"connection"`
	Host            string `mapstructure:"host" json:"host"`
	Port            int32  `mapstructure:"port" json:"port"`
	Database        string `mapstructure:"database" json:"database"`
	Username        string `mapstructure:"username" json:"username"`
	Password        string `mapstructure:"password" json:"password"`
	ForwardPassword string `mapstructure:"forward_password" json:"forward_password"`
	LogLevel        string
}

type SD struct {
	Driver string `mapstructure:"driver" json:"driver"`
}

type DC struct {
	Driver string `mapstructure:"driver" json:"driver"`
	Type   string `mapstructure:"type" json:"type"`
}

type Consul struct {
	Host    string `mapstructure:"host" json:"host"`
	Port    int32  `mapstructure:"port" json:"port"`
	Service struct {
		Name  string   `mapstructure:"name" json:"name"`
		Tags  []string `mapstructure:"tags" json:"tags"`
		Check struct {
			Protocol        string `mapstructure:"protocol" json:"protocol"`
			Endpoint        string `mapstructure:"endpoint" json:"endpoint"`
			Interval        int32  `mapstructure:"interval" json:"interval"`
			Timeout         int32  `mapstructure:"timeout" json:"timeout"`
			DeregisterAfter int32  `mapstructure:"deregister_after" json:"deregister_after"`
		} `mapstructure:"check" json:"check"`
	} `mapstructure:"service" json:"service"`
}

type Nacos struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port uint64 `mapstructure:"port"`
	} `mapstructure:"server"`
	Client struct {
		Namespace string `mapstructure:"namespace"`
		Username  string `mapstructure:"username"`
		Password  string `mapstructure:"password"`
		DataID    string `mapstructure:"data_id"`
		Group     string `mapstructure:"group"`
	} `mapstructure:"client"`
	ConfigClient config_client.IConfigClient
}

type Config struct {
	App    App    `mapstructure:"app" json:"app"`
	Log    Log    `mapstructure:"log" json:"log"`
	DB     DB     `mapstructure:"db" json:"db"`
	SD     SD     `mapstructure:"registry" json:"registry"`
	DC     DC     `mapstructure:"dc" json:"dc"`
	Consul Consul `mapstructure:"consul" json:"consul"`
	Nacos  Nacos  `mapstructure:"nacos" json:"nacos"`
}

func (c *Config) Load(filePath, filename, fileType string) error {
	viper.SetConfigName(filename)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(filePath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) Watch() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(c); err != nil {
			fmt.Println("failed to update config:", err)
		}
		fmt.Printf("%+v\n", c)
	})
	viper.WatchConfig()
}
