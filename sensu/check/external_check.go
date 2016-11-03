package check

import (
	"bytes"
	"os/exec"
	"syscall"
	"time"

	stdCheck "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
)

type ExternalCheck struct {
	Command string
}

func (c *ExternalCheck) Execute() stdCheck.CheckOutput {
	t0 := time.Now()
	cmd := exec.Command("/bin/sh", "-c", c.Command)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		return stdCheck.CheckOutput{
			stdCheck.Error,
			err.Error(),
			time.Now().Sub(t0).Seconds(),
			t0.Unix(),
			0,
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
		status,
		out.String(),
		time.Now().Sub(t0).Seconds(),
		t0.Unix(),
		0,
	}
}
