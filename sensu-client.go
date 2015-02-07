package main

import (
	"flag"
	"github.com/upfluence/sensu-client-go/sensu"
	"github.com/upfluence/sensu-client-go/sensu/transport"
	"os"
)

var (
	flagset = flag.NewFlagSet("etcdenv", flag.ExitOnError)
	flags   = sensu.ConfigFlagSet{}
)

func init() {
	flagset.BoolVar(&flags.Verbose, "v", false, "Verbose mode")
	flagset.StringVar(&flags.ConfigFile, "c", "", "Config file path")
}

func main() {
	flagset.Parse(os.Args[1:])

	cfg := sensu.NewConfigFromFlagSet(&flags)

	t := transport.NewRabbitMQTransport(cfg)

	client := sensu.NewClient(t, cfg)

	client.Start()
}
