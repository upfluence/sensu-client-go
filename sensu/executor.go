package sensu

import (
	"errors"

	stdCheck "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/check"
)

type CheckResponse struct {
	Check  stdCheck.CheckOutput `json:"check"`
	Client string               `json:"client"`
}

func executeCheck(input *stdCheck.CheckRequest) (*stdCheck.CheckOutput, error) {
	var output stdCheck.CheckOutput

	if ch, ok := check.Store[input.Extension]; input.Extension != "" && ok {
		output = ch.Execute()
	} else if ch, ok := check.Store[input.Name]; ok {
		output = ch.Execute()
	} else if input.Command == "" {
		return nil, errors.New("Command key not filled")
	} else {
		output = (&check.ExternalCheck{Request: input}).Execute()
	}

	output.CheckRequest = input

	return &output, nil
}
