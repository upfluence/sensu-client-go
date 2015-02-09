package handler

import (
	"fmt"
	"github.com/upfluence/sensu-client-go/sensu/check"
)

func Ok(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{0, fmt.Sprintf("OK: %s", message)}
}

func Warning(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{1, fmt.Sprintf("WARNING: %s", message)}
}

func Error(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{2, fmt.Sprintf("ERROR: %s", message)}
}
