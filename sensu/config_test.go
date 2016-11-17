package sensu

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
	stdClient "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/client"
)

func validateStringParameter(
	actualRabbitMqUri string,
	expectedRabbitMqUri string,
	parameterName string,
	t *testing.T) {

	if actualRabbitMqUri != expectedRabbitMqUri {
		t.Errorf("Expected %s to be \"%s\" but got \"%s\" instead!",
			parameterName,
			expectedRabbitMqUri,
			actualRabbitMqUri,
		)
	}
}

func TestRabbitMQURIDefaultValue(t *testing.T) {
	validateStringParameter((&Config{}).RabbitMQURI(),
		"amqp://guest:guest@localhost:5672/%2f",
		"RabbitMQ URI",
		t,
	)
}

func TestRabbitMQURIFromEnvVar(t *testing.T) {
	expectedRabbitMqUri := "amqp://user:password@example.com:5672"

	os.Setenv("RABBITMQ_URI", expectedRabbitMqUri)

	validateStringParameter((&Config{}).RabbitMQURI(),
		expectedRabbitMqUri,
		"RabbitMQ URI",
		t,
	)
}

func TestRabbitMQURIFromConfig(t *testing.T) {
	expectedRabbitMqUri := "amqp://user:password@example.com:5672"

	config := Config{config: &configPayload{RabbitMQURI: &expectedRabbitMqUri}}

	validateStringParameter(config.RabbitMQURI(),
		expectedRabbitMqUri,
		"RabbitMQ URI",
		t,
	)
}

var expectedClient = &stdClient.Client{Name: "test_client",
	Address:       "10.0.0.42",
	Subscriptions: strings.Split("email,messenger", ","),
}

func validateClient(actualClient *stdClient.Client, t *testing.T) {
	validateStringParameter(
		actualClient.Name,
		expectedClient.Name,
		"client name",
		t,
	)

	validateStringParameter(
		actualClient.Address,
		expectedClient.Address,
		"client address",
		t,
	)

	if !reflect.DeepEqual(
		actualClient.Subscriptions,
		expectedClient.Subscriptions,
	) {

		t.Errorf("Expected client subscriptions to be \"%v\""+
			" but got \"%v\" instead!",
			expectedClient.Subscriptions,
			actualClient.Subscriptions,
		)
	}
}

func expectedClientFromConfig(t *testing.T) {
	config := Config{config: &configPayload{Client: expectedClient}}

	validateClient(config.Client(), t)
}

func expectedClientFromEnvVars(t *testing.T) {
	os.Setenv("SENSU_CLIENT_NAME", expectedClient.Name)
	os.Setenv("SENSU_CLIENT_ADDRESS", expectedClient.Address)
	os.Setenv(
		"SENSU_CLIENT_SUBSCRIPTIONS",
		strings.Join(expectedClient.Subscriptions, ","),
	)

	validateClient((&Config{}).Client(), t)
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
		t.Errorf("Expected check count to be \"%d\" but got \"%d\" instead!",
			expectedCheckCount,
			actualCheckCount,
		)
	}
}
