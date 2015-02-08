package sensu

import (
	"flag"
	"os"
)

func ExtractFlags() *ConfigFlagSet {
	flagset := flag.NewFlagSet("etcdenv", flag.ExitOnError)
	flags := &ConfigFlagSet{}

	flagset.BoolVar(&flags.Verbose, "v", false, "Verbose mode")
	flagset.StringVar(&flags.ConfigFile, "c", "", "Config file path")

	flagset.Parse(os.Args[1:])

	return flags
}
