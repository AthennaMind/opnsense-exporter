package options

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/prometheus/common/promlog"
)

var (
	logLevel = kingpin.Flag(
		"log.level",
		"Log level. One of: [debug, info, warn, error]").
		Default("info").
		String()
	logFormat = kingpin.Flag(
		"log.format",
		"Log format. One of: [logfmt, json]").
		Default("logfmt").
		String()
)

func Logger() (log.Logger, error) {

	promlogConfig := &promlog.Config{
		Level:  &promlog.AllowedLevel{},
		Format: &promlog.AllowedFormat{},
	}

	if err := promlogConfig.Level.Set(*logLevel); err != nil {
		return nil, err
	}

	if err := promlogConfig.Format.Set(*logFormat); err != nil {
		return nil, err
	}

	return promlog.New(promlogConfig), nil
}
