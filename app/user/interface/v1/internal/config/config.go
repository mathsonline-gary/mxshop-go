package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      *App      `mapstructure:"app"`
	Server   *Server   `mapstructure:"server"`
	Data     *Data     `mapstructure:"data"`
	Registry *Registry `mapstructure:"registry"`
}

func (c *Config) Load(filePath string) error {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	filename := strings.TrimSuffix(base, ext)
	ext = strings.TrimPrefix(ext, ".")
	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
