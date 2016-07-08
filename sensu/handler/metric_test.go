package handler

import (
	"fmt"
	"testing"
	"time"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
)

func TestRending(t *testing.T) {
	m := &Metric{}
	m.AddPoint(&Point{"foo", 0.1})
	m.AddPoint(&Point{"bar", 1.0})

	tstamp := time.Now().Unix()
	r := m.Render()

	if r.Status != check.Success {
		t.Errorf("Wrong exit code: %d", r.Status)
	}

	expectedOutput := fmt.Sprintf(
		"foo 0.100000 %d\nbar 1.000000 %d",
		tstamp,
		tstamp,
	)

	if r.Output != expectedOutput {
		t.Errorf("Wrong output: %s and it's %s", r.Output, expectedOutput)
	}
}
