package check

import (
	"testing"

	stdCheck "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
)

func TestEmptyCommand(t *testing.T) {
	r := (&ExternalCheck{}).Execute()

	if r.Status != stdCheck.Error {
		t.Errorf("The status is not failed, %d", r.Status)
	}

	if r.Duration <= 0.0 {
		t.Errorf("The duration is not positive: %f", r.Duration)
	}
}

func TestCorrectCommand(t *testing.T) {
	r := (&ExternalCheck{"ls"}).Execute()

	if r.Status != stdCheck.Success {
		t.Errorf("The status is not success, %d", r.Status)
	}

	if r.Duration <= 0.0 {
		t.Errorf("The duration is not positive: %f", r.Duration)
	}
}

func TestOtherExitCodeCommand(t *testing.T) {
	r := (&ExternalCheck{"lsi /fiz/fux"}).Execute()

	if r.Status != stdCheck.Error {
		t.Errorf("The status is not success, %d", r.Status)
	}

	if r.Duration <= 0.0 {
		t.Errorf("The duration is not positive: %f", r.Duration)
	}
}
