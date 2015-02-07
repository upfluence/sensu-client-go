package sensu

import (
	"log"
	"os"
	"os/signal"
)

const CurrentVersion string = "0.0.1"

type Client struct {
	Transport  Transport
	Processors []Processor
	Config     *Config
}

func NewClient(transport Transport, cfg *Config) *Client {
	processors := []Processor{}

	for _, s := range cfg.Subscriptions() {
		processors = append(processors, &Subscriber{s, nil})
	}

	client := Client{
		transport,
		append(processors, &KeepAlive{}),
		cfg,
	}

	for _, processor := range client.Processors {
		processor.SetClient(&client)
	}

	return &client
}

func (c *Client) Start() error {
	c.Transport.Connect()

	for _, processor := range c.Processors {
		go processor.Start()
	}

	sig := make(chan os.Signal)

	signal.Notify(sig, os.Kill, os.Interrupt)

	s := <-sig

	log.Printf("Signal %s received", s.String())

	return c.Transport.Close()
}
