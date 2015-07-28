package check

type ExitStatus uint8

const (
	Success ExitStatus = iota
	Warning
	Error
)

type CheckOutput struct {
	Status   ExitStatus
	Output   string
	Duration float64
	Executed int64
}

type Check interface {
	Execute() CheckOutput
}
