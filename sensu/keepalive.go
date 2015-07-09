package sensu

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
)

type KeepAlive struct {
	Client    *Client
	closeChan chan bool
}

func NewKeepAlive() *KeepAlive {
	return &KeepAlive{nil, make(chan bool)}
}

func (k *KeepAlive) SetClient(c *Client) error {
	k.Client = c

	return nil
}

func (k *KeepAlive) PublishKeepAlive() {
	log.Println("Publishing keepalive")

	payload := make(map[string]interface{})

	payload["timestamp"] = time.Now().Unix()
	payload["version"] = CurrentVersion
	payload["name"] = k.Client.Config.Name()
	payload["address"] = k.Client.Config.Address()
	payload["subscriptions"] = k.Client.Config.Subscriptions()

	p, err := json.Marshal(payload)

	if err != nil {
		log.Printf("something goes wrong : %s", err.Error())
	}

	log.Printf("Payload sent: %s", bytes.NewBuffer(p).String())

	err = k.Client.Transport.Publish("direct", "keepalives", "", p)

	if err != nil {
		log.Printf("something goes wrong : %s", err.Error())
	}

}

func (k *KeepAlive) Start() error {
	t := time.Tick(20 * time.Second)

	k.PublishKeepAlive()

	for {
		select {
		case <-t:
			k.PublishKeepAlive()
		case <-k.closeChan:
			return nil
		}
	}

	return nil
}

func (k *KeepAlive) Close() {
	k.closeChan <- true
}
