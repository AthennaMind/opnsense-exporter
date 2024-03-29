package options

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
)

func Init() {
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
}
