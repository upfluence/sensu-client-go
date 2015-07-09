package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/upfluence/sensu-client-go/sensu/check"
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
	output := []string{}

	for _, p := range m.Points {
		output = append(output, p.Render())
	}

	return check.ExtensionCheckResult{0, strings.Join(output, "\n")}
}
