package behaviors

import (
	"fmt"
	"io"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	simple_behaviors "github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	opentracing "github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

var log = logger.GetLogger("flowmodel-opentracing")

// Flow implements model.FlowBehavior
type OpenTracingFlow struct {
}

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

// Start implements model.Flow.Start
func (fb *OpenTracingFlow) Start(ctx model.FlowContext) (started bool, taskEntries []*model.TaskEntry) {
	tracer, closer := initJaeger("flow-tracer")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	started, taskEntries = (&simple_behaviors.Flow{}).Start(ctx)

	sp := opentracing.StartSpan("flogo-flow")
	defer sp.Finish()

	sp.LogFields(opentracing_log.String("key", "value"))
	sp.LogKV("key", "value")

	ctx.WorkingData().AddAttr("opentracing-flow-span-context", data.TypeAny, sp.Context())

	return started, taskEntries
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
	(&simple_behaviors.Flow{}).Done(ctx)
}
