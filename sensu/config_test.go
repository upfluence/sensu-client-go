package sensu

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/goutils/testing/utils"
	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
	stdClient "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/client"
)

var dummyClient = &stdClient.Client{
	Name:          "test_client",
	Address:       "10.0.0.42",
	Subscriptions: strings.Split("email,messenger", ","),
}

func TestRabbitMQURIDefaultValue(t *testing.T) {
	utils.ValidateStringParameter(
		(&Config{}).RabbitMQURI(),
		defaultRabbitMQURI,
		"RabbitMQ URI",
		t,
	)
}

func TestRabbitMQURIFromEnvVar(t *testing.T) {
	expectedRabbitMqUri := "amqp://user:password@example.com:5672"

	os.Setenv("RABBITMQ_URI", expectedRabbitMqUri)
	defer os.Unsetenv("RABBITMQ_URI")

	utils.ValidateStringParameter(
		(&Config{}).RabbitMQURI(),
		expectedRabbitMqUri,
		"RabbitMQ URI",
		t,
	)
}

func TestRabbitMQURIFromConfig(t *testing.T) {
	expectedRabbitMqUri := "amqp://user:password@example.com:5672"

	config := Config{config: &configPayload{RabbitMQURI: &expectedRabbitMqUri}}

	utils.ValidateStringParameter(
		config.RabbitMQURI(),
		expectedRabbitMqUri,
		"RabbitMQ URI",
		t,
	)
}

func validateClient(actualClient *stdClient.Client, expectedClient *stdClient.Client, t *testing.T) {
	utils.ValidateStringParameter(
		actualClient.Name,
		expectedClient.Name,
		"client name",
		t,
	)

	utils.ValidateStringParameter(
		actualClient.Address,
		expectedClient.Address,
		"client address",
		t,
	)

	if !reflect.DeepEqual(
		actualClient.Subscriptions,
		expectedClient.Subscriptions,
	) {
		t.Errorf(
			"Expected client subscriptions to be \"%#v\" but got \"%#v\" instead!",
			expectedClient.Subscriptions,
			actualClient.Subscriptions,
		)
	}
}

func TestClientFromConfig(t *testing.T) {
	config := Config{config: &configPayload{Client: dummyClient}}

	validateClient(config.Client(), dummyClient, t)
}

func TestClientFromEnvVars(t *testing.T) {
	os.Setenv("SENSU_CLIENT_NAME", dummyClient.Name)
	defer os.Unsetenv("SENSU_CLIENT_NAME")

	os.Setenv("SENSU_CLIENT_ADDRESS", dummyClient.Address)
	defer os.Unsetenv("SENSU_CLIENT_ADDRESS")

	os.Setenv(
		"SENSU_CLIENT_SUBSCRIPTIONS",
		strings.Join(dummyClient.Subscriptions, ","),
	)
	defer os.Unsetenv("SENSU_CLIENT_SUBSCRIPTIONS")

	validateClient((&Config{}).Client(), dummyClient, t)
}

func TestClientFromEnvVarsNoSubscriptions(t *testing.T) {
	dummyClientNoSubscriptions := dummyClient
	dummyClientNoSubscriptions.Subscriptions = []string{}

	os.Setenv("SENSU_CLIENT_NAME", dummyClientNoSubscriptions.Name)
	defer os.Unsetenv("SENSU_CLIENT_NAME")

	os.Setenv("SENSU_CLIENT_ADDRESS", dummyClientNoSubscriptions.Address)
	defer os.Unsetenv("SENSU_CLIENT_ADDRESS")

	validateClient((&Config{}).Client(), dummyClientNoSubscriptions, t)
}

func TestChecksFromConfig(t *testing.T) {
	expectedCheckCount := 2
	config := Config{
		config: &configPayload{
			Checks: []*check.Check{&check.Check{}, &check.Check{}},
		},
	}

	actualCheckCount := len(config.Checks())

	if expectedCheckCount != actualCheckCount {
		t.Errorf(
			"Expected check count to be %d but got %d instead!",
			expectedCheckCount,
			actualCheckCount,
		)
	}
}

func TestNewConfigFromFile(t *testing.T) {
	if c, err := NewConfigFromFile(nil, ""); c != nil || err != errNoClientName {
		t.Errorf("Expected (nil, %v) but got (%v, %v)", errNoClientName, c, err)
	}
}

func TestRabbitMQHAConfigDefaultValue(t *testing.T) {
	haConfig, err := (&Config{}).RabbitMQHAConfig()

	if err != nil {
		t.Errorf(
			"Expected a nil error but got \"%s\" instead!",
			err,
		)
	}

	expectedConfigCont := 1

	if len(haConfig) != expectedConfigCont {
		t.Errorf(
			"Expected the config count to be %d but got %d instead!",
			expectedConfigCont,
			len(haConfig),
		)
	}

	utils.ValidateStringParameter(
		haConfig[0].GetURI(),
		defaultRabbitMQURI,
		"RabbitMQ URI",
		t,
	)
}
