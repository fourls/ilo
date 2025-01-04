package display

import (
	"fmt"
	"log"
	"time"

	"github.com/fourls/ilo/internal/ilofile"
)

type CliObserver struct {
	logger    *log.Logger
	project   *ilofile.Definition
	flow      *ilofile.Flow
	step      ilofile.Step
	flowStart time.Time
}

func NewObserver(project *ilofile.Definition, logger *log.Logger) CliObserver {
	return CliObserver{project: project, logger: logger}
}

func (o *CliObserver) FlowEntered(f *ilofile.Flow) {
	o.flow = f
	flowId := fmt.Sprintf("%s / %s", o.project.Name, o.flow.Name)
	HorizontalRule{Header: flowId}.Print(o.logger)

	o.flowStart = time.Now()
}

func (o *CliObserver) StepEntered(s ilofile.Step) {
	o.step = s
}

func (o *CliObserver) StepOutput(text string) {
	o.logger.Println(text)
}

func (o *CliObserver) StepPassed() {
	o.step = nil
}

func (o *CliObserver) StepFailed(err error) {
	o.logger.Println(err.Error())
}

func (o *CliObserver) FlowPassed() {
	o.flow = nil

	duration := time.Since(o.flowStart).Round(time.Millisecond)
	status := fmt.Sprintf("PASSED in %s", duration)
	HorizontalRule{Footer: status}.Print(o.logger)
}

func (o *CliObserver) FlowFailed() {
	o.flow = nil

	duration := time.Since(o.flowStart).Round(time.Millisecond)
	status := fmt.Sprintf("FAILED after %s", duration)
	HorizontalRule{Footer: status}.Print(o.logger)
}
