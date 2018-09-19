package behaviors

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	simple_behaviors "github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/opentracing/opentracing-go"
)

type OpenTracingTask struct {
}

// Enter implements model.Task.Enter
func (tb *OpenTracingTask) Enter(ctx model.TaskContext) (enterResult model.EnterResult) {
	// retrieve parent span context
	parentSpanContextAttr, _ := ctx.FlowWorkingData().GetAttr("opentracing-flow-span-context")
	parentSpanContext := parentSpanContextAttr.Value().(opentracing.SpanContext)

	// create span for task
	sp := opentracing.StartSpan(ctx.Task().Name(), opentracing.ChildOf(parentSpanContext))
	//sp.SetTag("tag", "value")

	// store span in working data to close it later
	spanAttr, err := data.NewAttribute("opentracing-task-span", data.TypeAny, sp)
	if err == nil {
		ctx.AddWorkingData(spanAttr)
	}

	// delegate to simple model
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
	taskSpanAttr, exists := ctx.GetWorkingData("opentracing-task-span")
	if exists {
		taskSpan := taskSpanAttr.Value().(opentracing.Span)
		taskSpan.Finish()
	}

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
