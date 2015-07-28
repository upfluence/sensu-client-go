package sensu

import (
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	CurrentVersion     string        = "0.1.0-dev"
	CONNECTION_TIMEOUT time.Duration = 5 * time.Second
)

type Client struct {
	Transport Transport
	Config    *Config
}

func NewClient(transport Transport, cfg *Config) *Client {

	client := Client{
		transport,
		cfg,
	}

	return &client
}

func (c *Client) buildProcessors() []Processor {
	processors := []Processor{NewKeepAlive()}

	for _, s := range c.Config.Subscriptions() {
		processors = append(processors, NewSubscriber(s))
	}

	for _, processor := range processors {
		processor.SetClient(c)
	}

	return processors
}

func (c *Client) Start() error {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	for {
		c.Transport.Connect()

		for !c.Transport.IsConnected() {
			time.Sleep(CONNECTION_TIMEOUT)
			c.Transport.Connect()
		}

		processors := c.buildProcessors()
		for _, processor := range processors {
			go processor.Start()
		}

		select {
		case s := <-sig:
			log.Printf("Signal %s received", s.String())

			for _, processor := range processors {
				processor.Close()
			}

			return c.Transport.Close()
		case <-c.Transport.GetClosingChan():
			log.Println("Transport disconnected")

			for _, processor := range processors {
				processor.Close()
			}

			c.Transport.Close()
		}
	}
}
