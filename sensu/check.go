package sensu

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Check struct {
	Payload map[string]interface{}
}

type CheckOutput struct {
	Status   int
	Output   string
	Duration float64
	Executed int64
}

func (c *Check) Execute() CheckOutput {
	command := strings.Split(c.Payload["command"].(string), " ")

	t0 := time.Now()
	cmd := exec.Command(command[0], command[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		return CheckOutput{
			127,
			err.Error(), time.Now().Sub(t0).Seconds(),
			t0.Unix(),
		}
	}

	status := 0

	if err := cmd.Wait(); err != nil {
		status = 127
		if exiterr, ok := err.(*exec.ExitError); ok {
			if statusReturn, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				status = int(statusReturn)
			}
		}
	}

	return CheckOutput{
		status,
		out.String(),
		time.Now().Sub(t0).Seconds(),
		t0.Unix(),
	}
}
