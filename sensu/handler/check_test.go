package handler

import (
	"testing"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
)

func TestOk(t *testing.T) {
	r := Ok("foo")

	if r.Status != check.Success {
		t.Errorf("Wrong status code")
	}

	if r.Output != "OK: foo" {
		t.Errorf("Wrong output message")
	}
}

func TestWarning(t *testing.T) {
	r := Warning("foo")

	if r.Status != check.Warning {
		t.Errorf("Wrong status code")
	}

	if r.Output != "WARNING: foo" {
		t.Errorf("Wrong output message")
	}
}

func TestError(t *testing.T) {
	r := Error("foo")

	if r.Status != check.Error {
		t.Errorf("Wrong status code")
	}

	if r.Output != "ERROR: foo" {
		t.Errorf("Wrong output message")
	}
}
