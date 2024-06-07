package config

type App struct {
	Env   string `mapstructure:"env"`
	Debug bool   `mapstructure:"debug"`
}
