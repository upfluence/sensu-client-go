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
}

func (s *Subscriber) SetClient(c *Client) error {
	s.Client = c

	return nil
}

func (s *Subscriber) Start() error {
	var output check.CheckOutput
	var b []byte

	funnel := strings.Join(
		[]string{
			s.Client.Config.Name(),
			CurrentVersion,
			strconv.Itoa(int(time.Now().Unix())),
		},
		"-",
	)

	msgChan := make(chan []byte)
	stopChan := make(chan bool)

	go s.Client.Transport.Subscribe("#", s.Subscription, funnel, msgChan, stopChan)

	log.Printf("Subscribed to %s", s.Subscription)

	failures := 0

	for {
		if failures >= MAX_FAILS {
			stopChan <- true
			msgChan = make(chan []byte)
			failures = 0
			s.Client.Transport.Close()
			time.Sleep(MAX_TIME)
			s.Client.Transport.Connect()
			go s.Client.Transport.Subscribe(
				"#",
				s.Subscription,
				funnel,
				msgChan,
				stopChan,
			)

		}

		b = <-msgChan

		payload := make(map[string]interface{})

		log.Printf("Check received : %s", bytes.NewBuffer(b).String())
		json.Unmarshal(b, &payload)

		if _, ok := payload["name"]; !ok {
			log.Printf("The name field is not filled")
			failures++
			continue
		}

		if ch, ok := check.Store[payload["name"].(string)]; ok {
			output = ch.Execute()
		} else if _, ok := payload["command"]; !ok {
			log.Printf("The command field is not filled")
			continue
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

	return nil
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
