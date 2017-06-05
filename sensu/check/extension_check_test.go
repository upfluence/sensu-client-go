package check

import (
	"testing"

	stdCheck "github.com/upfluence/sensu-go/sensu/check"
)

func FunctionT() ExtensionCheckResult {
	return ExtensionCheckResult{stdCheck.Warning, "foo"}
}

func TestExtension(t *testing.T) {
	r := (&ExtensionCheck{FunctionT}).Execute()

	if r.Status != stdCheck.Warning {
		t.Errorf("Wrong status: %d", r.Status)
	}

	if r.Output != "foo" {
		t.Errorf("Wrong output: %s", r.Output)
	}

	if r.Duration < 0.0 {
		t.Errorf("The duration is not positive: %f", r.Duration)
	}
}
