package behaviors

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"

	simple_behaviors "github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
)

type OpenTracingTask struct {
}

// Enter implements model.Task.Enter
func (tb *OpenTracingTask) Enter(ctx model.TaskContext) (enterResult model.EnterResult) {
	return (&simple_behaviors.Task{}).Enter(ctx)
}

// Eval implements model.Task.Eval
func (tb *OpenTracingTask) Eval(ctx model.TaskContext) (evalResult model.EvalResult, err error) {
	return (&simple_behaviors.Task{}).Eval(ctx)
}

// PostEval implements model.Task.PostEval
func (tb *OpenTracingTask) PostEval(ctx model.TaskContext) (evalResult model.EvalResult, err error) {
	return (&simple_behaviors.Task{}).PostEval(ctx)
}

// Done implements model.Task.Done
func (tb *OpenTracingTask) Done(ctx model.TaskContext) (notifyFlow bool, taskEntries []*model.TaskEntry, err error) {
	return (&simple_behaviors.Task{}).Done(ctx)
}

// Done implements model.Task.Skip
func (tb *OpenTracingTask) Skip(ctx model.TaskContext) (notifyFlow bool, taskEntries []*model.TaskEntry) {
	return (&simple_behaviors.Task{}).Skip(ctx)
}

// Done implements model.Task.Error
func (tb *OpenTracingTask) Error(ctx model.TaskContext, err error) (handled bool, taskEntries []*model.TaskEntry) {
	return (&simple_behaviors.Task{}).Error(ctx, err)
}
