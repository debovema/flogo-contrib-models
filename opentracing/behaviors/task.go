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
	parentSpanAttr, exists := ctx.FlowWorkingData().GetAttr("opentracing-flow-span")
	if exists {
		parentSpan := parentSpanAttr.Value().(opentracing.Span)

		var span opentracing.Span
		var spanName string

		// create child span for task (might be an activity or an iterator)
		if ctx.Task().TypeID() == "iterator" {
			span = opentracing.StartSpan(ctx.Task().Name()+" (iterator)", opentracing.ChildOf(parentSpan.Context()))
			span.SetTag("type", "flogo:iterator")
			spanName = "opentracing-iterator-span"
			span.SetTag("id", ctx.Task().ID()+" (iterator)")
		} else {
			span = opentracing.StartSpan(ctx.Task().Name(), opentracing.ChildOf(parentSpan.Context()))
			span.SetTag("type", "flogo:activity")
			spanName = "opentracing-activity-span"
			span.SetTag("id", ctx.Task().ID())
		}

		// store child span in working data to close it later
		spanAttr, err := data.NewAttribute(spanName, data.TypeAny, span)
		if err == nil {
			ctx.AddWorkingData(spanAttr)
		}
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
	taskSpanAttr, exists := ctx.GetWorkingData("opentracing-activity-span")
	if exists {
		taskSpan := taskSpanAttr.Value().(opentracing.Span)
		taskSpan.Finish()
	}

	// delegate to simple model
	return (&simple_behaviors.Task{}).PostEval(ctx)
}

// Done implements model.Task.Done
func (tb *OpenTracingTask) Done(ctx model.TaskContext) (notifyFlow bool, taskEntries []*model.TaskEntry, err error) {
	taskSpanAttr, exists := ctx.GetWorkingData("opentracing-activity-span")
	if exists {
		taskSpan := taskSpanAttr.Value().(opentracing.Span)
		taskSpan.Finish()
	}

	iterationSpanAttr, exists := ctx.GetWorkingData("opentracing-iteration-span")
	if exists {
		iterationSpan := iterationSpanAttr.Value().(opentracing.Span)
		iterationSpan.Finish()
	}

	iteratorSpanAttr, exists := ctx.GetWorkingData("opentracing-iterator-span")
	if exists {
		iteratorSpan := iteratorSpanAttr.Value().(opentracing.Span)
		iteratorSpan.Finish()
	}

	// delegate to simple model
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
