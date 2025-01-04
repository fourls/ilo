package ilofile

type Step interface {
	StepType() StepType
	String() string
}

type RunFlowStep interface {
	Args() []string
}

type EchoFlowStep interface {
	Message() string
}

type StepType int

const (
	StepRunProgram StepType = iota
	StepEchoMessage
)

type Flow struct {
	Name    string
	Dir     string
	Steps   []Step
	Project *Definition
}

type Definition struct {
	Name  string
	Path  string
	Flows map[string]Flow
}
