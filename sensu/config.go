package sensu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/upfluence/sensu-go/sensu/check"
	"github.com/upfluence/sensu-go/sensu/client"
	"github.com/upfluence/sensu-go/sensu/transport/rabbitmq"
)

const defaultRabbitMQURI string = "amqp://guest:guest@localhost:5672/%2F"

var errNoClientName = errors.New("No client name provided")

type configFlagSet struct {
	configFile string
	verbose    bool
}

type Config struct {
	flagSet *configFlagSet
	config  *configPayload
}

type configPayload struct {
	Client            *client.Client              `json:"client,omitempty"`
	Checks            []*check.Check              `json:"checks,omitempty"`
	RabbitMQURI       *string                     `json:"rabbitmq_uri,omitempty"`
	RabbitMQTransport []*rabbitmq.TransportConfig `json:"rabbitmq,omitempty"`
}

func fetchEnv(envs ...string) string {
	for _, env := range envs {
		if v := os.Getenv(env); v != "" {
			return v
		}
	}

	return ""
}

// split is a wrapper for strings.Split, which returns an empty array
// for empty string inputs instead of an array containing an empty string
func split(str string, token string) []string {
	if len(str) == 0 {
		return []string{}
	}

	return strings.Split(str, token)
}

func NewConfigFromFile(
	flagset *configFlagSet,
	configFile string,
) (*Config, error) {
	cfg := Config{flagset, &configPayload{}}

	if configFile != "" {
		buf, err := ioutil.ReadFile(configFile)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(buf, &cfg.config); err != nil {
			return nil, err
		}

		cfg.addDefaultSubscription()
	}

	if cfg.Client().Name == "" {
		return nil, errNoClientName
	}

	return &cfg, nil
}

func NewConfigFromFlagSet(flagset *configFlagSet) (*Config, error) {
	var file string

	if flagset != nil {
		file = flagset.configFile
	}

	return NewConfigFromFile(flagset, file)
}

func (c *Config) RabbitMQURI() string {
	if cfg := c.config; cfg != nil && cfg.RabbitMQURI != nil {
		return *cfg.RabbitMQURI
	} else if uri := fetchEnv("RABBITMQ_URI", "RABBITMQ_URL"); uri != "" {
		return uri
	}

	return defaultRabbitMQURI
}

// RabbitMQHAConfig superseeds RabbitMQURI() by first checking
// for a HA cluster configuration and then calling RabbitMQURI()
// if it can't find one
func (c *Config) RabbitMQHAConfig() ([]*rabbitmq.TransportConfig, error) {
	if cfg := c.config; cfg != nil && cfg.RabbitMQTransport != nil {
		return cfg.RabbitMQTransport, nil
	}

	config, err := rabbitmq.NewTransportConfig(c.RabbitMQURI())

	if err != nil {
		return []*rabbitmq.TransportConfig{}, err
	}

	return []*rabbitmq.TransportConfig{config}, nil
}

func (c *Config) Client() *client.Client {
	if c.config != nil {
		if c.config.Client != nil {
			return c.config.Client
		}
	} else {
		c.config = &configPayload{}
	}

	// Initialize the config from environment variables
	c.config.Client = &client.Client{
		Name:          os.Getenv("SENSU_CLIENT_NAME"),
		Address:       fetchEnv("SENSU_CLIENT_ADDRESS", "SENSU_ADDRESS"),
		Subscriptions: split(os.Getenv("SENSU_CLIENT_SUBSCRIPTIONS"), ","),
	}

	c.addDefaultSubscription()

	return c.config.Client
}

func (c *Config) Checks() []*check.Check {
	if cfg := c.config; cfg != nil {
		return cfg.Checks
	}

	return []*check.Check{}
}

// addDefaultSubscription emulates ruby client behavior:
// add default subscription - client:name
// Without at least one subscription sensu server will crash.
func (c *Config) addDefaultSubscription() {
	c.config.Client.Subscriptions = removeDuplicates(
		append(
			c.config.Client.Subscriptions,
			fmt.Sprintf("client:%s", c.config.Client.Name),
		),
	)
}

// Removes duplicates from string slices
func removeDuplicates(xs []string) []string {

	temp := make(map[string]bool)
	result := []string{}

	for _, x := range xs {
		temp[x] = true
	}

	for k, _ := range temp {
		result = append(result, k)
	}

	return result
}
