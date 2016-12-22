package rabbitmq

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/streadway/amqp"
	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/goutils/log"
)

type RabbitMQTransport struct {
	Connection     *amqp.Connection
	Channel        *amqp.Channel
	ClosingChannel chan bool
	Configs        []*TransportConfig
}

func NewRabbitMQTransport(uri string) (*RabbitMQTransport, error) {
	config, err := NewTransportConfig(uri)

	if err != nil {
		return nil, fmt.Errorf("Received invalid URI: %s", err)
	}

	return NewRabbitMQHATransport([]*TransportConfig{config}), nil
}

func NewRabbitMQHATransport(configs []*TransportConfig) *RabbitMQTransport {
	return &RabbitMQTransport{
		ClosingChannel: make(chan bool),
		Configs:        configs,
	}
}

func (t *RabbitMQTransport) GetClosingChan() chan bool {
	return t.ClosingChannel
}

func (t *RabbitMQTransport) Connect() error {
	var (
		uri           string
		err           error
		randGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for _, idx := range randGenerator.Perm(len(t.Configs)) {
		config := t.Configs[idx]
		uri = config.GetURI()

		log.Noticef("Trying to connect to URI: %s", uri)

		// TODO: Figure out how to specify the Prefetch value as well
		// See amqp.Channel.Qos (it doesn't seem to be used currently)

		// TODO: Add SSL support via amqp.DialTLS

		heartbeatString := config.Heartbeat.String()
		if heartbeatString != "" {
			var heartbeat time.Duration
			heartbeat, err = time.ParseDuration(heartbeatString + "s")

			if err != nil {
				log.Warningf("Failed to parse the heartbeat: %s", uri, err.Error())
				continue
			}

			t.Connection, err = amqp.DialConfig(
				uri,
				amqp.Config{Heartbeat: heartbeat},
			)
		} else {
			// Use amqp.defaultHeartbeat (10s)
			t.Connection, err = amqp.Dial(uri)
		}

		if err != nil {
			log.Warningf("Failed to connect to URI %s: %s", uri, err.Error())
			continue
		}

		break
	}

	if err != nil {
		log.Errorf("RabbitMQ connection error: %s", err.Error())
		return err
	}

	t.Channel, err = t.Connection.Channel()

	if err != nil {
		log.Errorf("RabbitMQ channel error: %s", err.Error())
		return err
	}

	log.Noticef("RabbitMQ connection and channel opened to %s", uri)

	closeChan := make(chan *amqp.Error)
	t.Channel.NotifyClose(closeChan)

	go func() {
		<-closeChan
		t.ClosingChannel <- true
	}()

	return nil
}

func (t *RabbitMQTransport) IsConnected() bool {
	if t.Connection == nil || t.Channel == nil {
		return false
	}

	return true
}

func (t *RabbitMQTransport) Close() error {
	if t.Connection == nil {
		return errors.New("The connection is not opened")
	}

	defer func() {
		t.Channel = nil
		t.Connection = nil
	}()
	t.Connection.Close()
	return nil
}

func (t *RabbitMQTransport) Publish(exchangeType, exchangeName, key string, message []byte) error {
	if t.Channel == nil {
		return errors.New("The channel is not opened")
	}

	if err := t.Channel.ExchangeDeclare(exchangeName, exchangeType, false, false, false, false, nil); err != nil {
		return err
	}

	err := t.Channel.Publish(exchangeName, key, false, false, amqp.Publishing{Body: message})

	return err
}

func (t *RabbitMQTransport) Subscribe(key, exchangeName, queueName string, messageChan chan []byte, stopChan chan bool) error {
	if t.Channel == nil {
		return errors.New("The channel is not opened")
	}

	if err := t.Channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		false,
		false,
		false,
		false,
		amqp.Table{},
	); err != nil {
		log.Errorf("Can't declare the exchange: %s", err.Error())
		return err
	}

	log.Infof("Exchange %s declared", exchangeName)

	if _, err := t.Channel.QueueDeclare(
		queueName,
		false,
		true,
		false,
		false,
		nil,
	); err != nil {
		log.Errorf("Can't declare the queue: %s", err.Error())
		return err
	}

	log.Infof("Queue %s declared", queueName)

	if err := t.Channel.QueueBind(queueName, key, exchangeName, false, nil); err != nil {
		log.Errorf("Can't bind the queue: %s", err.Error())
		return err
	}

	log.Noticef("Queue %s binded to %s for key %s", queueName, exchangeName, key)

	deliveryChange, err := t.Channel.Consume(queueName, "", true, false, false, false, nil)

	log.Infof("Consuming the queue %s", queueName)

	if err != nil {
		log.Errorf("Can't consume the queue: %s", err.Error())
		return err
	}

	for {
		select {
		case delivery, ok := <-deliveryChange:
			if ok {
				messageChan <- delivery.Body
			} else {
				t.ClosingChannel <- true
				break
			}
		case <-stopChan:
			break
		}
	}
}
