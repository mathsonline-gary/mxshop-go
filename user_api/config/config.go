package config

type UserSvcConfig struct {
	Name string `mapstructure:"name"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name"`
	UserSvcConfig UserSvcConfig `mapstructure:"user_svc"`
	ConsulConfig ConsulConfig `mapstructure:"consul" json:"consul_config"`
}
