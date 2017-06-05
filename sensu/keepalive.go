package sensu

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/upfluence/goutils/log"
	"github.com/upfluence/sensu-go/sensu/client"
)

const defaultInterval = 20 * time.Second

type KeepAlive struct {
	Client    *Client
	closeChan chan bool
}

type keepAlivePayload struct {
	*client.Client
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
}

func NewKeepAlive(c *Client) *KeepAlive {
	return &KeepAlive{c, make(chan bool)}
}

func (k *KeepAlive) publishKeepAlive() {
	log.Info("Publishing keepalive")

	p, err := json.Marshal(
		keepAlivePayload{
			k.Client.Config.Client(),
			time.Now().Unix(),
			currentVersion,
		},
	)

	if err != nil {
		log.Warningf("Something went wrong: %s", err.Error())
		return
	}

	err = k.Client.Transport.Publish("direct", "keepalives", "", p)
	log.Infof("Payload sent: %s", bytes.NewBuffer(p).String())

	if err != nil {
		log.Warningf("Something went wrong: %s", err.Error())
	}
}

func (k *KeepAlive) Start() error {
	t := time.Tick(defaultInterval)

	k.publishKeepAlive()

	for {
		select {
		case <-t:
			k.publishKeepAlive()
		case <-k.closeChan:
			return nil
		}
	}
}

func (k *KeepAlive) Close() {
	k.closeChan <- true
}
