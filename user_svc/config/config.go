package config

type NacosServerConfig struct {
	Host string `mapstructure:"host"`
	Port uint64 `mapstructure:"port"`
}

type NacosClientConfig struct {
	Namespace string `mapstructure:"namespace"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	DataID    string `mapstructure:"data_id"`
	Group     string `mapstructure:"group"`
}

type NacosConfig struct {
	NacosServerConfig `mapstructure:"server"`
	NacosClientConfig `mapstructure:"client"`
}

type DBConfig struct {
	Driver   string `mapstructure:"driver" json:"driver"`
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Database string `mapstructure:"database" json:"database"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
}

type AppConfig struct {
	Name  string `mapstructure:"name" json:"name"`
	Host  string `mapstructure:"host" json:"host"`
	Port  int    `mapstructure:"port" json:"port"`
	Env   string `mapstructure:"env" json:"env"`
	Debug bool   `mapstructure:"debug" json:"debug"`
}

type ConsulConfig struct {
	Host    string `mapstructure:"host" json:"host"`
	Port    int    `mapstructure:"port" json:"port"`
	UserSvc struct {
		Check struct {
			Host string `mapstructure:"host" json:"host"`
		} `mapstructure:"check" json:"check"`
	} `mapstructure:"user_svc" json:"user_svc"`
}

type Config struct {
	DBConfig     DBConfig     `mapstructure:"db" json:"db"`
	AppConfig    AppConfig    `mapstructure:"app" json:"app"`
	ConsulConfig ConsulConfig `mapstructure:"consul" json:"consul"`
}
