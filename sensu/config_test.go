package sensu

import (
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/upfluence/goutils/testing/utils"
	"github.com/upfluence/sensu-go/sensu/check"
	stdClient "github.com/upfluence/sensu-go/sensu/client"
)

type dummyFlagSet struct{}

func (*dummyFlagSet) BoolVar(_ *bool, _ string, _ bool, _ string)       {}
func (*dummyFlagSet) StringVar(_ *string, _ string, _ string, _ string) {}
func (*dummyFlagSet) Parse([]string) error {
	return nil
}

func newDummyClient() *stdClient.Client {
	return &stdClient.Client{
		Name:          "test_client",
		Address:       "10.0.0.42",
		Subscriptions: strings.Split("email,messenger", ","),
	}
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

func validateClient(
	actualClient *stdClient.Client,
	expectedClient *stdClient.Client,
	t *testing.T,
) {
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

	// Sort subscription slices first, because DeepEqual also requires
	// the order of slices to be equal.
	sort.Strings(actualClient.Subscriptions)
	sort.Strings(expectedClient.Subscriptions)
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
	dummyClient := newDummyClient()

	config := Config{config: &configPayload{Client: dummyClient}}

	validateClient(config.Client(), dummyClient, t)
}

func TestClientFromEnvVars(t *testing.T) {
	dummyClient := newDummyClient()
	dummyClient.Subscriptions = append(
		dummyClient.Subscriptions,
		"client:test_client",
	)

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
	dummyClient := newDummyClient()
	dummyClient.Subscriptions = []string{"client:test_client"}

	os.Setenv("SENSU_CLIENT_NAME", dummyClient.Name)
	defer os.Unsetenv("SENSU_CLIENT_NAME")

	os.Setenv("SENSU_CLIENT_ADDRESS", dummyClient.Address)
	defer os.Unsetenv("SENSU_CLIENT_ADDRESS")

	validateClient((&Config{}).Client(), dummyClient, t)
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

func TestNewConfigFromFlagSet(t *testing.T) {
	origCmdlineParser := cmdlineParser
	cmdlineParser = &dummyFlagSet{}
	defer func() { cmdlineParser = origCmdlineParser }()

	dummyClientName := "test_client"

	os.Setenv("SENSU_CLIENT_NAME", dummyClientName)
	defer os.Unsetenv("SENSU_CLIENT_NAME")

	cfg, err := NewConfigFromFlagSet(ExtractFlags())

	if err != nil {
		t.Errorf("Expected error to be nil but got \"%s\" instead", err)
	}

	utils.ValidateStringParameter(
		cfg.config.Client.Name,
		dummyClientName,
		"client name",
		t,
	)
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

func TestSubscriptionBehaviour(t *testing.T) {
	for _, tCase := range []struct {
		in  string
		out []string
	}{
		{
			"testdata/client-noSubs.json",
			[]string{"client:foo"},
		},
		{
			"testdata/client-dupeSubs.json",
			strings.Split("unique,duplicate,client:foo", ","),
		},
		{
			"testdata/client-uniqueSubs.json",
			strings.Split("unique1,unique2,unique3,client:foo", ","),
		},
	} {

		tCaseCfg, err := NewConfigFromFile(nil, tCase.in)

		if err != nil {
			t.Errorf(
				"Expected a nil error but got \"%s\" instead!",
				err,
			)
		}

		validateClientSubscriptions(
			tCaseCfg.config.Client.Subscriptions,
			tCase.out,
			t,
		)
	}
}

func validateClientSubscriptions(s1 []string, s2 []string, t *testing.T) {
	sort.Strings(s1)
	sort.Strings(s2)

	if !reflect.DeepEqual(s1, s2) {

		t.Errorf(
			"Expected client subscriptions to be \"%#v\" but got \"%#v\" instead!",
			s1,
			s2,
		)
	}
}
