package main

import (
	"github.com/upfluence/sensu-client-go/sensu"
	"github.com/upfluence/sensu-go/sensu/transport/rabbitmq"
)

func main() {
	cfg := sensu.NewConfigFromFlagSet(sensu.ExtractFlags())

	t := rabbitmq.NewRabbitMQTransport(cfg.RabbitMQURI())
	client := sensu.NewClient(t, cfg)

	client.Start()
}
