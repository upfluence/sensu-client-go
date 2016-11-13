package sensu

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/goutils/log"
	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
)

type Standalone struct {
	check     *check.Check
	client    *Client
	closeChan chan bool
}

func NewStandalone(check *check.Check, c *Client) *Standalone {
	return &Standalone{check, c, make(chan bool)}
}

func (s *Standalone) Start() error {
	t := time.Tick(defaultInterval)

	if s.check.Interval > 0 {
		t = time.Tick(time.Duration(s.check.Interval) * time.Second)
	}

	log.Noticef("Setup standalone check %s", s.check.Name)

	for {
		select {
		case <-t:
			s.execute()
		case <-s.closeChan:
			log.Warningf("Graceful stop of %s", s.check.Name)
			return nil
		}
	}
}

func (s *Standalone) Close() {
	s.closeChan <- true
}

func (s *Standalone) execute() {
	if p, err := json.Marshal(s.check); err == nil {
		log.Infof("Check received: %s", bytes.NewBuffer(p).String())
	}

	output, err := executeCheck(
		&check.CheckRequest{Check: s.check, Issued: time.Now().Unix()},
	)

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
