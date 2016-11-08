package sensu

import (
	"flag"
	"os"
)

func ExtractFlags() *configFlagSet {
	flagset := flag.NewFlagSet("sensu-client-go", flag.ExitOnError)
	flags := &configFlagSet{}

	flagset.BoolVar(&flags.verbose, "v", false, "Verbose mode")
	flagset.StringVar(&flags.configFile, "c", "", "Config file path")

	flagset.Parse(os.Args[1:])

	return flags
}
