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
		Status:   output.Status,
		Output:   output.Output,
		Duration: time.Since(t0).Seconds(),
		Executed: t0.Unix(),
	}
}
