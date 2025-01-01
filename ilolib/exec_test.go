package ilolib

import (
	"reflect"
	"testing"
)

type MockTestExecutor struct {
	ExecuteFunc func(ExecParams) error
}

func (e MockTestExecutor) StepExecute(params ExecParams) error {
	return e.ExecuteFunc(params)
}

func TestRunStep_PassToExecutor(t *testing.T) {
	flow := Flow{
		Steps: []FlowStep{
			flowStep{stepType: StepEchoMessage},
		},
	}

	params := ExecParams{
		Env:       []string{"foo=bar"},
		Directory: "baz",
		Observer:  nil,
		Toolbox:   Toolbox{},
	}

	executorFactory := buildMockExecutorFactory(
		t,
		StepEchoMessage,
		params,
	)

	err := FlowExecutor{}.runStep(flow, 0, params, executorFactory)
	if err != nil {
		t.Fatalf("runStep got: %v, want: nil", err)
	}
}

func buildMockExecutorFactory(t *testing.T, expectedType StepType, expectedParams ExecParams) func(step FlowStep) StepExecutor {
	return func(step FlowStep) StepExecutor {
		if step.StepType() != expectedType {
			t.Fatalf("got: step with type %d, want: step with type %d", step.StepType(), expectedType)
		}

		return MockTestExecutor{
			ExecuteFunc: func(params ExecParams) error {
				if !reflect.DeepEqual(expectedParams, params) {
					t.Fatalf("got: %v, want: %v", params, expectedParams)
				}

				return nil
			},
		}
	}
}
