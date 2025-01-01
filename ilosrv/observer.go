package ilosrv

import (
	"log/slog"

	"fourls.dev/ilo/ilolib"
)

type StructuredObserver struct {
	logger    *slog.Logger
	project   *ilolib.ProjectDefinition
	flow      *ilolib.Flow
	step      ilolib.FlowStep
	stepIndex int
}

func newObserver(project *ilolib.ProjectDefinition, logger *slog.Logger) StructuredObserver {
	return StructuredObserver{project: project, logger: logger.With("project", project.Path), stepIndex: -1}
}

func (o *StructuredObserver) FlowEntered(f *ilolib.Flow) {
	o.flow = f
	o.stepIndex = -1
	o.logger.Info("Flow entered", "flow", o.flow.Name)
}

func (o *StructuredObserver) StepEntered(s ilolib.FlowStep) {
	o.step = s
	o.stepIndex += 1
	o.logger.Info("Step entered", "flow", o.flow.Name, "step", o.stepIndex, "stepText", o.step.String())
}

func (o *StructuredObserver) StepOutput(text string) {
	o.logger.Debug("> "+text, "flow", o.flow.Name, "step", o.stepIndex)
}

func (o *StructuredObserver) StepPassed() {
	o.logger.Info("Step passed", "flow", o.flow.Name, "step", o.stepIndex)
	o.step = nil
}

func (o *StructuredObserver) StepFailed(err error) {
	o.logger.Info("Step failed", "flow", o.flow.Name, "step", o.stepIndex, "error", err)
	o.step = nil
}

func (o *StructuredObserver) FlowPassed() {
	o.logger.Info("Flow passed", "flow", o.flow.Name)
	o.flow = nil
}

func (o *StructuredObserver) FlowFailed() {
	o.logger.Info("Flow failed", "flow", o.flow.Name, "step", o.stepIndex)
	o.flow = nil
}
