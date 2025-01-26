package options

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog/flag"
)

func Init() {
	flag.AddFlags(kingpin.CommandLine, PromLogConfig)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
}
