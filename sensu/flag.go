package sensu

import (
	"flag"
	"os"
)

var (
	cmdlineParser flagSet
)

type flagSet interface {
	BoolVar(p *bool, name string, value bool, usage string)
	StringVar(p *string, name string, value string, usage string)
	Parse(arguments []string) error
}

func init() {
	cmdlineParser = flag.NewFlagSet("sensu-client-go", flag.ExitOnError)
}

func ExtractFlags() *configFlagSet {
	flags := &configFlagSet{}

	cmdlineParser.BoolVar(&flags.verbose, "v", false, "Verbose mode")
	cmdlineParser.StringVar(&flags.configFile, "c", "", "Config file path")

	_ = cmdlineParser.Parse(os.Args[1:])

	return flags
}
