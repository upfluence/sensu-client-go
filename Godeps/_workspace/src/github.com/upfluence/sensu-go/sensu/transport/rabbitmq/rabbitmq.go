package rabbitmq

import (
	"errors"
	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/goutils/log"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/streadway/amqp"
)

type RabbitMQTransport struct {
	URI            string
	Connection     *amqp.Connection
	Channel        *amqp.Channel
	ClosingChannel chan bool
}

func NewRabbitMQTransport(uri string) *RabbitMQTransport {
	return &RabbitMQTransport{URI: uri, ClosingChannel: make(chan bool)}
}

func (t *RabbitMQTransport) GetClosingChan() chan bool {
	return t.ClosingChannel
}

func (t *RabbitMQTransport) Connect() error {
	var err error

	t.Connection, err = amqp.Dial(t.URI)

	if err != nil {
		log.Errorf("RabbitMQ connection error: %s", err.Error())
		return err
	}

	t.Channel, err = t.Connection.Channel()

	if err != nil {
		log.Errorf("RabbitMQ channel error: %s", err.Error())
		return err
	}

	log.Noticef("RabbitMQ connection and channel opened to %s", t.URI)

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
