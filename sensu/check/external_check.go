package check

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	"time"

	stdCheck "github.com/upfluence/sensu-go/sensu/check"
)

type ExternalCheck struct {
	Command string
}

func (c *ExternalCheck) Execute() stdCheck.CheckOutput {
	command := strings.Split(c.Command, " ")

	t0 := time.Now()
	cmd := exec.Command(command[0], command[1:]...)
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
				status = stdCheck.ExitStatus(statusReturn)
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
