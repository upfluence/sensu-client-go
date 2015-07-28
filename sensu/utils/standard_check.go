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
	Value func() float64
	// Used to render the output message of the check
	CheckMessage func(float64) string
}

func (c *StandardCheck) Check() check.ExtensionCheckResult {
	v := c.Value()

	if v > c.ErrorThreshold {
		return handler.Error(c.CheckMessage(v))
	} else if v > c.WarningThreshold {
		return handler.Warning(c.CheckMessage(v))
	} else {
		return handler.Ok(c.CheckMessage(v))
	}
}

func (c *StandardCheck) Metric() check.ExtensionCheckResult {
	m := &handler.Metric{}

	m.AddPoint(&handler.Point{c.MetricName, c.Value()})

	return m.Render()
}
