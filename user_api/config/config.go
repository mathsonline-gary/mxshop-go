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

type AppConfig struct {
	Name  string `json:"name"`
	Env   string `json:"env"`
	Debug bool   `json:"debug"`
}

type UserSvcConfig struct {
	Name string `json:"name"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Config struct {
	AppConfig     AppConfig     `json:"app"`
	UserSvcConfig UserSvcConfig `json:"user_svc"`
	ConsulConfig  ConsulConfig  `json:"consul"`
}
