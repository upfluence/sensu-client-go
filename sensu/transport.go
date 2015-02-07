package sensu

type Transport interface {
	Connect() error
	IsConnected() (bool, error)
	Close() error
	Publish(exchangeType, exchangeName, key string, message []byte) error
	Subscribe(key, exchangeName, queueName string, messageChan chan []byte, stopChan chan bool) error
}
