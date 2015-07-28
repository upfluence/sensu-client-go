package handler

import (
	"fmt"
	"github.com/upfluence/sensu-client-go/sensu/check"
)

func Ok(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{
		check.Success,
		fmt.Sprintf("OK: %s", message),
	}
}

func Warning(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{
		check.Warning,
		fmt.Sprintf("WARNING: %s", message),
	}
}

func Error(message string) check.ExtensionCheckResult {
	return check.ExtensionCheckResult{
		check.Error,
		fmt.Sprintf("ERROR: %s", message),
	}
}
