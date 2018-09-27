package behaviors

import (
	"github.com/debovema/flogo-contrib-models/opentracing/utils"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	simple_behaviors "github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/opentracing/opentracing-go"
)

var log = logger.GetLogger("flowmodel-opentracing")
var GlobalTracer opentracing.Tracer
var GlobalConfig *utils.OpenTracingConfig

// OpenTracingFlow implements model.FlowBehavior
type OpenTracingFlow struct {
}

// Start implements model.Flow.Start
func (fb *OpenTracingFlow) Start(ctx model.FlowContext) (started bool, taskEntries []*model.TaskEntry) {
	var currentTracer opentracing.Tracer
	var opentracingConfig *utils.OpenTracingConfig

	if GlobalTracer != nil {
		currentTracer = GlobalTracer
	} else {
		if GlobalConfig != nil {
			opentracingConfig = GlobalConfig
		} else {
			opentracingConfig = utils.ReadOpentracingContext(ctx)
		}

		if opentracingConfig != nil {
			tracer, err := utils.InitTracer(ctx.FlowDefinition().Name(), opentracingConfig)
			if err != nil || tracer == nil {
				log.Warn("Unable to init OpenTracing tracer. Ignoring.")
			} else {
				currentTracer = tracer
			}
		} else {
			log.Warn("Unable to init OpenTracing tracer. Ignoring.")
		}
	}

	if currentTracer != nil {
		opentracing.SetGlobalTracer(currentTracer)

		span := opentracing.StartSpan(ctx.FlowDefinition().Name())
		span.SetTag("type", "flogo:flow")

		// store span in working data to close it later and to pass the span context to activities
		ctx.WorkingData().AddAttr("opentracing-flow-span", data.TypeAny, span)
	}

	return (&simple_behaviors.Flow{}).Start(ctx)
}

// StartErrorHandler implements model.Flow.StartErrorHandler
func (fb *OpenTracingFlow) StartErrorHandler(ctx model.FlowContext) (taskEntries []*model.TaskEntry) {
	return (&simple_behaviors.Flow{}).StartErrorHandler(ctx)
}

// Resume implements model.Flow.Resume
func (fb *OpenTracingFlow) Resume(ctx model.FlowContext) (resumed bool) {
	return (&simple_behaviors.Flow{}).Resume(ctx)
}

// TasksDone implements model.Flow.TasksDone
func (fb *OpenTracingFlow) TaskDone(ctx model.FlowContext) (flowDone bool) {
	return (&simple_behaviors.Flow{}).TaskDone(ctx)
}

// Done implements model.Flow.Done
func (fb *OpenTracingFlow) Done(ctx model.FlowContext) {
	flowSpanAttr, exists := ctx.WorkingData().GetAttr("opentracing-flow-span")
	if exists {
		flowSpan := flowSpanAttr.Value().(opentracing.Span)
		flowSpan.Finish()
	}

	(&simple_behaviors.Flow{}).Done(ctx)
}
