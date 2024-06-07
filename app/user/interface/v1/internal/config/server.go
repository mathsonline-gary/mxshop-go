package config

type Server struct {
	HTTP ServerHTTP `mapstructure:"HTTP"`
}

type ServerHTTP struct {
	Host string `mapstructure:"host"`
	Port uint32 `mapstructure:"port"`
}
