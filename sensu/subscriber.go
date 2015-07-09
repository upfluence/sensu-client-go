package sensu

import (
	"bytes"
	"encoding/json"
	"github.com/upfluence/sensu-client-go/sensu/check"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_FAILS = 100
	MAX_TIME  = 60 * time.Second
)

type Subscriber struct {
	Subscription string
	Client       *Client
	closeChan    chan bool
}

func NewSubscriber(subscription string) *Subscriber {
	return &Subscriber{subscription, nil, make(chan bool)}
}

func (s *Subscriber) SetClient(c *Client) error {
	s.Client = c

	return nil
}

func (s *Subscriber) Start() error {
	funnel := strings.Join(
		[]string{
			s.Client.Config.Name(),
			CurrentVersion,
			strconv.Itoa(int(time.Now().Unix())),
		},
		"-",
	)

	msgChan, stopChan := s.subscribe(funnel)
	log.Printf("Subscribed to %s", s.Subscription)

	for {
		select {
		case b := <-msgChan:
			s.handleMessage(b)
		case <-s.closeChan:
			log.Printf("Gracefull stop of %s", s.Subscription)
			stopChan <- true
			return nil
		}
	}

	return nil
}

func (s *Subscriber) Close() {
	s.closeChan <- true
}

func (s *Subscriber) handleMessage(blob []byte) {
	var output check.CheckOutput
	payload := make(map[string]interface{})

	log.Printf("Check received : %s", bytes.NewBuffer(blob).String())
	json.Unmarshal(blob, &payload)

	if _, ok := payload["name"]; !ok {
		log.Printf("The name field is not filled")
		return
	}

	if ch, ok := check.Store[payload["name"].(string)]; ok {
		output = ch.Execute()
	} else if _, ok := payload["command"]; !ok {
		log.Printf("The command field is not filled")
		return
	} else {
		output = (&check.ExternalCheck{payload["command"].(string)}).Execute()
	}

	p, err := json.Marshal(s.forgeCheckResponse(payload, &output))

	if err != nil {
		log.Printf("something goes wrong : %s", err.Error())
	} else {
		log.Printf("Payload sent: %s", bytes.NewBuffer(p).String())

		err = s.Client.Transport.Publish("direct", "results", "", p)
	}
}

func (s *Subscriber) subscribe(funnel string) (chan []byte, chan bool) {
	msgChan := make(chan []byte)
	stopChan := make(chan bool)

	go func() {
		for {
			s.Client.Transport.Subscribe(
				"#",
				s.Subscription,
				funnel,
				msgChan,
				stopChan,
			)
		}
	}()

	return msgChan, stopChan
}

func (s *Subscriber) forgeCheckResponse(payload map[string]interface{}, output *check.CheckOutput) map[string]interface{} {
	result := make(map[string]interface{})

	result["client"] = s.Client.Config.Name()

	formattedOuput := make(map[string]interface{})

	formattedOuput["name"] = payload["name"]
	formattedOuput["issued"] = int(payload["issued"].(float64))
	formattedOuput["output"] = output.Output
	formattedOuput["duration"] = output.Duration
	formattedOuput["status"] = output.Status
	formattedOuput["executed"] = output.Executed

	result["check"] = formattedOuput

	return result
}
