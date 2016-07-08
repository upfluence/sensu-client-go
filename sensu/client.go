package sensu

import (
	"os"
	"os/signal"
	"time"

	"github.com/upfluence/goutils/log"
	"github.com/upfluence/sensu-go/sensu/transport"
)

const (
	CurrentVersion     string        = "0.1.2"
	CONNECTION_TIMEOUT time.Duration = 5 * time.Second
)

type Client struct {
	Transport transport.Transport
	Config    *Config
}

func NewClient(transport transport.Transport, cfg *Config) *Client {
	client := Client{
		transport,
		cfg,
	}

	return &client
}

func (c *Client) buildProcessors() []Processor {
	processors := []Processor{NewKeepAlive(c)}

	for _, s := range c.Config.Subscriptions() {
		processors = append(processors, NewSubscriber(s, c))
	}

	return processors
}

func (c *Client) Start() error {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	for {
		c.Transport.Connect()

		for !c.Transport.IsConnected() {
			select {
			case <-time.After(CONNECTION_TIMEOUT):
				c.Transport.Connect()
			case <-sig:
				return c.Transport.Close()
			}
		}

		processors := c.buildProcessors()
		for _, processor := range processors {
			go processor.Start()
		}

		select {
		case s := <-sig:
			log.Noticef("Signal %s received", s.String())

			for _, processor := range processors {
				processor.Close()
			}

			return c.Transport.Close()
		case <-c.Transport.GetClosingChan():
			log.Notice("Transport disconnected")

			for _, processor := range processors {
				processor.Close()
			}

			c.Transport.Close()
		}
	}
}
