package sensu

import (
	"testing"

	stdCheck "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/handler"
)

func validateCheckOutput(
	checkRequest *stdCheck.CheckRequest,
	expectedOutput *stdCheck.CheckOutput,
	t *testing.T) {

	output, err := executeCheck(checkRequest)

	if err != nil {
		t.Fatalf("Expected error to be nil but got \"%s\" instead!", err)
	}

	if output == nil {
		t.Fatalf("Expected non-nil output but got nil instead!")
	}

	if output.Status != expectedOutput.Status {
		t.Errorf("Expected status to be \"%d\" but got \"%d\" instead!",
			expectedOutput.Status,
			output.Status)
	}

	if output.Output != expectedOutput.Output {
		t.Errorf("Expected output to be \"%s\" but got \"%s\" instead!",
			expectedOutput.Output,
			output.Output)
	}
}

var checkTestFunction = func() check.ExtensionCheckResult {
	return handler.Ok("Test")
}

var extensionCheckTestData = []struct {
	extensionCheck *check.ExtensionCheck
	check          *stdCheck.Check
	output         string
}{
	{
		&check.ExtensionCheck{Function: checkTestFunction},
		&stdCheck.Check{Name: "extension_check"},
		"OK: Test"},
	{
		&check.ExtensionCheck{Function: checkTestFunction},
		&stdCheck.Check{Name: "extension_check",
			Extension: "named_extension_check"},
		"OK: Test"},
}

func TestExecuteExtensionCheck(t *testing.T) {
	for _, test := range extensionCheckTestData {
		check.Store[test.check.Name] = test.extensionCheck

		validateCheckOutput(
			&stdCheck.CheckRequest{Check: test.check},
			&stdCheck.CheckOutput{Status: 0, Output: test.output},
			t)
	}
}

func TestExecuteExternalCheck(t *testing.T) {
	validateCheckOutput(
		&stdCheck.CheckRequest{Check: &stdCheck.Check{Command: "printf Test"}},
		&stdCheck.CheckOutput{Status: 0, Output: "Test"},
		t)
}

func TestExecuteEmptyCheck(t *testing.T) {
	output, err := executeCheck(
		&stdCheck.CheckRequest{Check: &stdCheck.Check{}, Issued: 1479057736})

	if err != commandKeyError {
		t.Fatalf("Expected error to be \"%s\" but got \"%s\" instead!",
			commandKeyError,
			err)
	}

	if output != nil {
		t.Fatalf("Expected output to be nil but got \"%+v\" instead!", output)
	}
}
