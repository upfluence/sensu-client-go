package check

import (
	"bytes"
	"os/exec"
	"syscall"
	"time"

	stdCheck "github.com/upfluence/sensu-go/sensu/check"
)

type ExternalCheck struct {
	Request *stdCheck.CheckRequest
}

func (c *ExternalCheck) Execute() stdCheck.CheckOutput {
	t0 := time.Now()
	cmd := exec.Command("/bin/sh", "-c", c.Request.Command)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		return stdCheck.CheckOutput{
			CheckRequest: c.Request,
			Status:       stdCheck.Error,
			Output:       err.Error(),
			Duration:     time.Since(t0).Seconds(),
			Executed:     t0.Unix(),
		}
	}

	status := stdCheck.Success

	if err := cmd.Wait(); err != nil {
		status = stdCheck.Error
		if exiterr, ok := err.(*exec.ExitError); ok {
			if statusReturn, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				status = stdCheck.ExitStatus(statusReturn.ExitStatus())
			}
		}
	}

	return stdCheck.CheckOutput{
		Status:   status,
		Output:   out.String(),
		Duration: time.Since(t0).Seconds(),
		Executed: t0.Unix(),
	}
}
