package config

type DBConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	DBConfig  DBConfig  `mapstructure:"db" json:"db_config"`
	AppConfig AppConfig `mapstructure:"app" json:"app_config"`
	ConsulConfig ConsulConfig `mapstructure:"consul" json:"consul_config"`
}
