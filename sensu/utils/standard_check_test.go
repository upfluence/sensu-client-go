package utils

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/upfluence/sensu-client-go/sensu/check"
)

var (
	meth1    = func() float64 { return 0.9 }
	meth2    = func() float64 { return 1.9 }
	meth3    = func() float64 { return 2.9 }
	messMeth = func(f float64) string {
		return fmt.Sprintf("messMeth: %.2f", f)
	}

	okCheck      = &StandardCheck{2.0, 1.0, "ok", meth1, messMeth}
	warningCheck = &StandardCheck{2.0, 1.0, "warn", meth2, messMeth}
	errCheck     = &StandardCheck{2.0, 1.0, "err", meth3, messMeth}
)

func TestEnvVarEmpty(t *testing.T) {
	if v := EnvironmentValueOrConst("FOOO", 42.0); v != 42.0 {
		t.Errorf("Wrong value returned %f", v)
	}
}

func TestEnvVarMalformed(t *testing.T) {
	if v := EnvironmentValueOrConst("HOSTNAME", 42.0); v != 42.0 {
		t.Errorf("Wrong value returned %f", v)
	}
}

func TestEnvVarOk(t *testing.T) {
	os.Setenv("FUZ", "28")

	if v := EnvironmentValueOrConst("HOSTNAME", 28.0); v != 28.0 {
		t.Errorf("Wrong value returned %f", v)
	}
}

func TestStandardCheckOk(t *testing.T) {
	v := okCheck.Check()

	if v.Status != check.Success {
		t.Errorf("Wrong status: %d", v.Status)
	}

	if v.Output != "OK: messMeth: 0.90" {
		t.Errorf("Wrong message: %s", v.Output)
	}
}

func TestStandardCheckWarning(t *testing.T) {
	v := warningCheck.Check()

	if v.Status != check.Warning {
		t.Errorf("Wrong status: %d", v.Status)
	}

	if v.Output != "WARNING: messMeth: 1.90" {
		t.Errorf("Wrong message: %s", v.Output)
	}
}

func TestStandardCheckError(t *testing.T) {
	v := errCheck.Check()

	if v.Status != check.Error {
		t.Errorf("Wrong status: %d", v.Status)
	}

	if v.Output != "ERROR: messMeth: 2.90" {
		t.Errorf("Wrong message: %s", v.Output)
	}
}

func TestStandardCheckMtric(t *testing.T) {
	v := errCheck.Metric()

	if v.Status != check.Success {
		t.Errorf("Wrong status: %d", v.Status)
	}

	if v.Output != fmt.Sprintf("err 2.900000 %d", time.Now().Unix()) {
		t.Errorf("Wrong message: %s", v.Output)
	}
}
