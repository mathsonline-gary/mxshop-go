package config

type UserSvcConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name"`
	UserSvcConfig UserSvcConfig `mapstructure:"user_svc"`
}
