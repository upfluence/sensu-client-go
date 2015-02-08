package check

type CheckOutput struct {
	Status   int
	Output   string
	Duration float64
	Executed int64
}

type Check interface {
	Execute() CheckOutput
}
