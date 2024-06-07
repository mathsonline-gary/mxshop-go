package config

type Data struct {
	UserService DataUserService `mapstructure:"user-service"`
}

type DataUserService struct {
	Driver string `mapstructure:"driver"`
	Host   string `mapstructure:"host"`
	Port   uint32 `mapstructure:"port"`
	Name   string `mapstructure:"name"`
}
