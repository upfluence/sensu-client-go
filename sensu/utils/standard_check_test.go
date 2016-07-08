package utils

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
)

var (
	meth1    = func() (float64, error) { return 0.9, nil }
	meth2    = func() (float64, error) { return 1.9, nil }
	meth3    = func() (float64, error) { return 2.9, nil }
	failMeth = func() (float64, error) { return 0.0, errors.New("Foo Bar") }
	asc      = func(x, y float64) bool { return x < y }
	desc     = func(x, y float64) bool { return x > y }

	messMeth = func(f float64) string {
		return fmt.Sprintf("messMeth: %.2f", f)
	}

	okCheck        = &StandardCheck{2.0, 1.0, "ok", meth1, messMeth, asc}
	warningCheck   = &StandardCheck{2.0, 1.0, "warn", meth2, messMeth, asc}
	errCheck       = &StandardCheck{2.0, 1.0, "err", meth3, messMeth, asc}
	errInvertCheck = &StandardCheck{2.0, 1.0, "err", meth1, messMeth, desc}
	failCheck      = &StandardCheck{2.0, 1.0, "ok", failMeth, messMeth, asc}
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

	if v := EnvironmentValueOrConst("FUZ", 21.0); v != 28.0 {
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

func TestStandardCheckInvertError(t *testing.T) {
	v := errInvertCheck.Check()

	if v.Status != check.Error {
		t.Errorf("Wrong status: %d", v.Status)
	}

	if v.Output != "ERROR: messMeth: 0.90" {
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

func TestStandardCheckMtricFail(t *testing.T) {
	v := failCheck.Metric()

	if v.Status != check.Success {
		t.Errorf("Wrong status: %d", v.Status)
	}

	if v.Output != "" {
		t.Errorf("Wrong message: %s", v.Output)
	}
}

func TestStandardCheckFail(t *testing.T) {
	v := failCheck.Check()

	if v.Status != check.Error {
		t.Errorf("Wrong status: %d", v.Status)
	}

	if v.Output != "ERROR: Foo Bar" {
		t.Errorf("Wrong message: %s", v.Output)
	}
}
