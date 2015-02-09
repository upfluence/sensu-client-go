package handler

import (
	"fmt"
	"github.com/upfluence/sensu-client-go/sensu/check"
	"time"
)

type Point struct {
	Name  string
	Value float64
}

type Metric struct {
	Points []Point
}

func (p Point) Render() string {
	return fmt.Sprintf("%s %f %d", p.Name, p.Value, time.Now().Unix())
}

func (m *Metric) AddPoint(p Point) {
	m.Points = append(m.Points, p)
}

func (m *Metric) Render() check.ExtensionCheckResult {
	output := ""

	for _, p := range m.Points {
		output = fmt.Sprintf("%s\n%s", output, p.Render())
	}

	return check.ExtensionCheckResult{0, output}
}
