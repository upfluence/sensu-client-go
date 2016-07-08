package check

import (
	"time"

	stdCheck "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
)

type ExtensionCheckResult struct {
	Status stdCheck.ExitStatus
	Output string
}

type ExtensionCheck struct {
	Function func() ExtensionCheckResult
}

func (c *ExtensionCheck) Execute() stdCheck.CheckOutput {
	t0 := time.Now()

	output := c.Function()

	return stdCheck.CheckOutput{
		output.Status,
		output.Output,
		time.Now().Sub(t0).Seconds(),
		t0.Unix(),
		0,
	}
}
