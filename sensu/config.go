package sensu

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const DefaultRabbitMQURI string = "amqp://guest:guest@localhost:5672/%2f"

type ConfigFlagSet struct {
	ConfigFile string
	Verbose    bool
}

type Config struct {
	FlagSet *ConfigFlagSet
	config  map[string]interface{}
}

func NewConfigFromFlagSet(flagset *ConfigFlagSet) *Config {
	cfg := Config{
		flagset,
		make(map[string]interface{}),
	}

	if flagset != nil && flagset.ConfigFile != "" {
		buf, err := ioutil.ReadFile(flagset.ConfigFile)

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(buf, &cfg.config)

		if err != nil {
			panic(err)
		}

		cfg.config = cfg.config["client"].(map[string]interface{})
	}

	return &cfg
}

func (c *Config) RabbitMQURI() string {
	uri := os.Getenv("RABBITMQ_URL")

	if uri == "" && c.config["rabbit_uri"] != nil {
		uri = c.config["rabbit_uri"].(string)
	} else if uri == "" {
		uri = DefaultRabbitMQURI
	}

	return uri
}

func (c *Config) Name() string {
	name := os.Getenv("SENSU_CLIENT_NAME")

	if name == "" && c.config["name"] != nil {
		name = c.config["name"].(string)
	}

	return name
}

func (c *Config) Address() string {
	address := os.Getenv("SENSU_CLIENT_ADDRESS")

	if address == "" && c.config["address"] != nil {
		address = c.config["address"].(string)
	}

	return address
}

func (c *Config) Subscriptions() []string {
	subscriptions := []string{}

	for _, s := range strings.Split(os.Getenv("SENSU_CLIENT_SUBSCRIPTIONS"), ",") {
		if s != "" {
			subscriptions = append(subscriptions, s)
		}
	}

	if len(subscriptions) == 0 && c.config["subscriptions"] != nil {
		for _, s := range c.config["subscriptions"].([]interface{}) {
			subscriptions = append(subscriptions, s.(string))
		}
	}

	return subscriptions
}
