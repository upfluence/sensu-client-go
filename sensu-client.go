package main

import (
	"log"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/transport/rabbitmq"
	"github.com/upfluence/sensu-client-go/sensu"
)

func main() {
	cfg, err := sensu.NewConfigFromFlagSet(sensu.ExtractFlags())

	if err != nil {
		log.Fatal(err.Error())
	}

	t, err := rabbitmq.NewRabbitMQTransport(cfg.RabbitMQURI())

	if err != nil {
		log.Fatal(err.Error())
	}

	client := sensu.NewClient(t, cfg)

	client.Start()
}
