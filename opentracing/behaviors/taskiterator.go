package behaviors

import (
	"strconv"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	simple_behaviors "github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/opentracing/opentracing-go"
)

// OpenTracingIteratorTask implements model.TaskBehavior
type OpenTracingIteratorTask struct {
	OpenTracingTask
}

// Eval implements model.TaskBehavior.Eval
func (tb *OpenTracingIteratorTask) Eval(ctx model.TaskContext) (evalResult model.EvalResult, err error) {
	// retrieve parent span context (the one from iterator)
	iteratorSpanAttr, exists := ctx.GetWorkingData("opentracing-iterator-span")
	if exists {
		var iterationIndex int
		var span opentracing.Span
		var parentSpan opentracing.SpanReference

		iteratorSpan := iteratorSpanAttr.Value().(opentracing.Span)

		// check whether a span exists: the one from the last iteration which must be finnished here
		iterationSpanAttr, exists := ctx.GetWorkingData("opentracing-iteration-span")
		if exists {
			iterationSpan := iterationSpanAttr.Value().(opentracing.Span)
			iterationSpan.Finish()

			// retrieve iteration index
			iterationIndexAttr, exists := ctx.GetWorkingData("opentracing-iteration-index")
			if exists {
				iterationIndex = iterationIndexAttr.Value().(int) + 1
			}

			// create "follows from" span for next iterations
			//parentSpan = opentracing.FollowsFrom(iterationSpan.Context())
			// create child span for next iterations
			parentSpan = opentracing.ChildOf(iteratorSpan.Context())
		} else {
			// create child span for first iteration
			parentSpan = opentracing.ChildOf(iteratorSpan.Context())
		}

		span = opentracing.StartSpan(ctx.Task().Name()+" (iteration #"+strconv.Itoa(iterationIndex)+")", parentSpan)

		span.SetTag("type", "flogo:activity")
		span.SetTag("id", ctx.Task().ID()+"@"+strconv.Itoa(iterationIndex))

		// store child span in working data to close it later
		spanAttr, err := data.NewAttribute("opentracing-iteration-span", data.TypeAny, span)
		if err == nil {
			ctx.AddWorkingData(spanAttr)
		}
		indexAttr, err := data.NewAttribute("opentracing-iteration-index", data.TypeInteger, iterationIndex)
		if err == nil {
			ctx.AddWorkingData(indexAttr)
		}
	}

	return (&simple_behaviors.IteratorTask{}).Eval(ctx)
}

// PostEval implements model.TaskBehavior.PostEval
func (tb *OpenTracingIteratorTask) PostEval(ctx model.TaskContext) (evalResult model.EvalResult, err error) {
	iteratorSpanAttr, exists := ctx.GetWorkingData("opentracing-iterator-span")
	if exists {
		iteratorSpan := iteratorSpanAttr.Value().(opentracing.Span)
		iteratorSpan.Finish()
	}

	return (&simple_behaviors.IteratorTask{}).PostEval(ctx)
}
