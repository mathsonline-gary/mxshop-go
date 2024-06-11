package config

type Server struct {
	HTTP ServerHTTP `mapstructure:"HTTP"`
}

type ServerHTTP struct {
	Scheme string `mapstructure:"scheme"`
	Host   string `mapstructure:"host"`
	Port   uint32 `mapstructure:"port"`
}
