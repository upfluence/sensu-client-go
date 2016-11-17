package utils

import (
	"os"
	"strconv"

	"github.com/upfluence/sensu-client-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/handler"
)

func EnvironmentValueOrConst(envVar string, constVal float64) float64 {
	if v, err := strconv.Atoi(os.Getenv(envVar)); err == nil {
		return float64(v)
	}

	return constVal
}

type StandardCheck struct {
	// Minimal value to trigger an error
	ErrorThreshold float64
	// Minimal value to trigger a warning
	WarningThreshold float64
	// The metric name sent to the server
	MetricName string
	// Function to compute the value along the checks
	Value func() (float64, error)
	// Used to render the output message of the check
	CheckMessage func(float64) string
	// Use to evaluate the value with the thresholds
	Comp func(float64, float64) bool
}

func (c *StandardCheck) Check() check.ExtensionCheckResult {
	v, err := c.Value()

	if err != nil {
		return handler.Error(err.Error())
	}

	if c.Comp(c.ErrorThreshold, v) {
		return handler.Error(c.CheckMessage(v))
	} else if c.Comp(c.WarningThreshold, v) {
		return handler.Warning(c.CheckMessage(v))
	} else {
		return handler.Ok(c.CheckMessage(v))
	}
}

func (c *StandardCheck) Metric() check.ExtensionCheckResult {
	m := &handler.Metric{}

	v, err := c.Value()

	if err == nil {
		m.AddPoint(&handler.Point{Name: c.MetricName, Value: v})
	}
	return m.Render()
}
