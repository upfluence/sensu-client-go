package sensu

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
	stdClient "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/client"
)

var dummyClient = &stdClient.Client{
	Name:          "test_client",
	Address:       "10.0.0.42",
	Subscriptions: strings.Split("email,messenger", ","),
}

func validateStringParameter(
	actual string,
	expected string,
	parameterName string,
	t *testing.T) {

	if actual != expected {
		t.Errorf(
			"Expected %s to be \"%s\" but got \"%s\" instead!",
			parameterName,
			expected,
			actual,
		)
	}
}

func TestRabbitMQURIDefaultValue(t *testing.T) {
	validateStringParameter(
		(&Config{}).RabbitMQURI(),
		"amqp://guest:guest@localhost:5672/%2f",
		"RabbitMQ URI",
		t,
	)
}

func TestRabbitMQURIFromEnvVar(t *testing.T) {
	expectedRabbitMqUri := "amqp://user:password@example.com:5672"

	os.Setenv("RABBITMQ_URI", expectedRabbitMqUri)
	defer os.Unsetenv("RABBITMQ_URI")

	validateStringParameter(
		(&Config{}).RabbitMQURI(),
		expectedRabbitMqUri,
		"RabbitMQ URI",
		t,
	)
}

func TestRabbitMQURIFromConfig(t *testing.T) {
	expectedRabbitMqUri := "amqp://user:password@example.com:5672"

	config := Config{config: &configPayload{RabbitMQURI: &expectedRabbitMqUri}}

	validateStringParameter(
		config.RabbitMQURI(),
		expectedRabbitMqUri,
		"RabbitMQ URI",
		t,
	)
}

func validateClient(actualClient *stdClient.Client, expectedClient *stdClient.Client, t *testing.T) {
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
			"Expected check count to be \"%d\" but got \"%d\" instead!",
			expectedCheckCount,
			actualCheckCount,
		)
	}
}
