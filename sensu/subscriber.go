package sensu

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/goutils/log"
	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
)

const (
	maxFail = 100
	maxTime = 60 * time.Second
)

type Subscriber struct {
	subscription string
	client       *Client
	closeChan    chan bool
}

func NewSubscriber(subscription string, c *Client) *Subscriber {
	return &Subscriber{subscription, c, make(chan bool)}
}

func (s *Subscriber) Start() error {
	funnel := strings.Join(
		[]string{
			s.client.Config.Client().Name,
			currentVersion,
			strconv.Itoa(int(time.Now().Unix())),
		},
		"-",
	)

	msgChan, stopChan := s.subscribe(funnel)
	log.Noticef("Subscribed to %s", s.subscription)

	for {
		select {
		case b := <-msgChan:
			s.handleMessage(b)
		case <-s.closeChan:
			log.Warningf("Gracefull stop of %s", s.subscription)
			stopChan <- true
			return nil
		}
	}
}

func (s *Subscriber) Close() {
	s.closeChan <- true
}

func (s *Subscriber) handleMessage(blob []byte) {
	var input check.CheckRequest

	log.Noticef("Check received: %s", bytes.NewBuffer(blob).String())

	if err := json.Unmarshal(blob, &input); err != nil {
		log.Errorf("Something went wrong: %s", err.Error())
		return
	}

	output, err := executeCheck(&input)

	if err != nil {
		log.Error(err.Error())
		return
	}

	p, err := json.Marshal(
		&CheckResponse{Check: *output, Client: s.client.Config.Client().Name},
	)

	if err != nil {
		log.Errorf("Something went wrong: %s", err.Error())
	} else {
		log.Noticef("Payload sent: %s", bytes.NewBuffer(p).String())
		s.client.Transport.Publish("direct", "results", "", p)
	}
}

func (s *Subscriber) subscribe(funnel string) (chan []byte, chan bool) {
	msgChan := make(chan []byte)
	stopChan := make(chan bool)

	go func() {
		for {
			s.client.Transport.Subscribe(
				"#",
				s.subscription,
				funnel,
				msgChan,
				stopChan,
			)
		}
	}()

	return msgChan, stopChan
}
