package sensu

import (
	"testing"

	stdCheck "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
	stdClient "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/client"
)

func TestMissingCommandKey(t *testing.T) {
	standaloneProcessor := &Standalone{check: &stdCheck.Check{}}

	err := standaloneProcessor.execute()

	if err != commandKeyError {
		t.Fatalf(
			"Expected error to be \"%s\" but got \"%s\" instead!",
			commandKeyError,
			err,
		)
	}
}

type transportPublishParameters struct {
	exchangeType, exchangeName, key string
	message                         []byte
}

type dummyTransport struct {
	publishParameters *transportPublishParameters
}

func (t *dummyTransport) Connect() error {
	return nil
}

func (t *dummyTransport) IsConnected() bool {
	return true
}

func (t *dummyTransport) Close() error {
	return nil
}

func (t *dummyTransport) Publish(
	exchangeType,
	exchangeName,
	key string,
	message []byte) error {

	t.publishParameters = &transportPublishParameters{
		exchangeType,
		exchangeName,
		key,
		message,
	}

	return nil
}

func (t *dummyTransport) Subscribe(
	key,
	exchangeName,
	queueName string,
	messageChan chan []byte,
	stopChan chan bool) error {

	return nil
}

func (t *dummyTransport) GetClosingChan() chan bool {
	return nil
}

func TestRunCommand(t *testing.T) {
	standaloneProcessor := NewStandalone(
		&stdCheck.Check{Command: "ls"},
		&Client{
			Config: &Config{
				config: &configPayload{
					Client: &stdClient.Client{Name: "Test"},
				},
			},
			Transport: &dummyTransport{},
		},
	)

	err := standaloneProcessor.execute()

	if err != nil {
		t.Fatalf("Expected error to be nil but got \"%s\" instead!", err)
	}

	expectedPublishParams := transportPublishParameters{
		exchangeType: "direct",
		exchangeName: "results",
		key:          "",
	}

	publishParams :=
		standaloneProcessor.client.Transport.(*dummyTransport).publishParameters

	if publishParams.exchangeType != expectedPublishParams.exchangeType {
		t.Errorf(
			"Expected exchange type to be \"%s\" but got \"%s\" instead!",
			expectedPublishParams.exchangeType,
			publishParams.exchangeType,
		)
	}

	if publishParams.exchangeName != expectedPublishParams.exchangeName {
		t.Errorf(
			"Expected exchange name to be \"%s\" but got \"%s\" instead!",
			expectedPublishParams.exchangeName,
			publishParams.exchangeName,
		)
	}

	if publishParams.key != expectedPublishParams.key {
		t.Errorf(
			"Expected key to be \"%s\" but got \"%s\" instead!",
			expectedPublishParams.key,
			publishParams.key,
		)
	}

	// Not particularly relevant here to validate the contents of the message
	if publishParams.message == nil {
		t.Errorf("Expected message type to be initialized but got nil instead!")
	}
}
