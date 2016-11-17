package handler

import (
	"fmt"

	stdCheck "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/check"
)

func Ok(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{
		Status: stdCheck.Success,
		Output: fmt.Sprintf("OK: %s", message),
	}
}

func Warning(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{
		Status: stdCheck.Warning,
		Output: fmt.Sprintf("WARNING: %s", message),
	}
}

func Error(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{
		Status: stdCheck.Error,
		Output: fmt.Sprintf("ERROR: %s", message),
	}
}
