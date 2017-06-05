package check

import (
	"testing"

	stdCheck "github.com/upfluence/sensu-go/sensu/check"
)

func TestEmptyCommand(t *testing.T) {
	r := (&ExternalCheck{
		&stdCheck.CheckRequest{Check: &stdCheck.Check{}},
	}).Execute()

	if r.Status != stdCheck.Success {
		t.Errorf("The status is not failed, %d", r.Status)
	}

	if r.Duration < 0.0 {
		t.Errorf("The duration is not positive: %f", r.Duration)
	}
}

func TestCorrectCommand(t *testing.T) {
	r := (&ExternalCheck{
		&stdCheck.CheckRequest{Check: &stdCheck.Check{Command: "ls"}},
	}).Execute()

	if r.Status != stdCheck.Success {
		t.Errorf("The status is not success, %d", r.Status)
	}

	if r.Duration < 0.0 {
		t.Errorf("The duration is not positive: %f", r.Duration)
	}
}

func TestOtherExitCodeCommand(t *testing.T) {
	r := (&ExternalCheck{
		&stdCheck.CheckRequest{Check: &stdCheck.Check{Command: "lsi /fiz/fux"}},
	}).Execute()

	if r.Status != 127 {
		t.Errorf("The status is not success, %d", r.Status)
	}

	if r.Duration < 0.0 {
		t.Errorf("The duration is not positive: %f", r.Duration)
	}
}

func TestCustomeExitCodeCommand(t *testing.T) {
	r := (&ExternalCheck{
		&stdCheck.CheckRequest{
			Check: &stdCheck.Check{Command: `/bin/bash -c "echo 'foo'; exit 42"`},
		},
	}).Execute()

	if r.Status != 42 {
		t.Errorf("The status is not success, %d", r.Status)
	}

	if r.Output != "foo\n" {
		t.Errorf("Wrong output: %v", r.Output)
	}
}
