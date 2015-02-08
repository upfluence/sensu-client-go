package check

import "time"

type ExtensionCheckResult struct {
	Status int
	Output string
}

type ExtensionCheck struct {
	Function func() ExtensionCheckResult
}

func (c *ExtensionCheck) Execute() CheckOutput {
	t0 := time.Now()

	output := c.Function()

	return CheckOutput{
		output.Status,
		output.Output,
		time.Now().Sub(t0).Seconds(),
		t0.Unix(),
	}
}
