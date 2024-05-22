package log

import (
	"fmt"

	"github.com/spf13/pflag"
)

// Option is log option. It is used to initialize the logger. The values are parsed from the configuration file or command line arguments.
type Option struct {
	Level Level `json:"level" mapstructure:"level" toml:"level" yaml:"level"`
}

// Validate validates the log option.
func (o *Option) Validate() []error {
	var errs []error

	if o.Level < LevelDebug || o.Level > LevelFatal {
		errs = append(errs, fmt.Errorf("invalid log level %d", o.Level))
	}

	return errs
}

func (o *Option) ParseFlagSet(fs *pflag.FlagSet) {
	var lvl int8
	fs.Int8Var(&lvl, "log.level", int8(o.Level), "log level")
	o.Level = Level(lvl)
}
