package behaviors

import (
	"errors"
	"fmt"
	"io"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	simple_behaviors "github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var log = logger.GetLogger("flowmodel-opentracing")

// OpenTracingFlow implements model.FlowBehavior
type OpenTracingFlow struct {
}

const (
	hostPort = "0.0.0.0:0" // not applicable -> leave as-is

	// Debug mode.
	debug = false

	// same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
	sameSpan = true

	// make Tracer generate 128 bit traceID's for root spans.
	traceID128Bit = true
)

type OpenTracingConfig struct {
	Implementation string   `json:"implementation"`
	Transport      string   `json:"transport"`
	Endpoints      []string `json:"endpoints"`
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

func initZipkinHttp(serviceName string, endpoint string) opentracing.Tracer {
	// Create our HTTP collector.
	collector, err := zipkin.NewHTTPCollector(endpoint)
	if err != nil {
		panic(fmt.Sprintf("unable to create Zipkin HTTP collector: %+v\n", err))

	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)

	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to create Zipkin tracer: %+v\n", err))
	}

	return tracer
}

func initZipkinKafka(serviceName string, endpoint []string) opentracing.Tracer {
	// Create our Kafka collector.
	collector, err := zipkin.NewKafkaCollector(endpoint)
	//collector, err := zipkin.NewKafkaCollector(endpoint, zipkin.KafkaLogger(zipkin.LogWrapper(log.New(os.Stdout, log.Prefix(), log.Flags()))))

	if err == nil {
		// Create our recorder.
		recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)

		// Create our tracer.
		tracer, err := zipkin.NewTracer(
			recorder,
			zipkin.ClientServerSameSpan(sameSpan),
			zipkin.TraceID128Bit(traceID128Bit),
		)
		if err != nil {
			panic(fmt.Sprintf("unable to create Zipkin tracer: %+v\n", err))
		}
		return tracer
	} else {
		// panic(fmt.Sprintf("unable to create Zipkin Kafka collector: %+v\n", err))
		return nil
	}
}

func readOpentracingContext(ctx model.FlowContext) *OpenTracingConfig {
	opentracingConfigData, exists := ctx.FlowDefinition().GetAttr("opentracing-config")

	if exists {
		value := opentracingConfigData.Value()

		opentracingConfig := &OpenTracingConfig{}
		mapstructure.Decode(value, opentracingConfig)

		return opentracingConfig
	} else {
		return nil
	}
}

func initTracer(serviceName string, opentracingConfig *OpenTracingConfig) (opentracing.Tracer, error) {
	switch opentracingConfig.Implementation {
	case "zipkin":
		switch opentracingConfig.Transport {
		case "http":
			return initZipkinHttp(serviceName, opentracingConfig.Endpoints[0]), nil
		case "kafka":
			return initZipkinKafka(serviceName, opentracingConfig.Endpoints), nil
		default:
			return nil, errors.New("supported transports for OpenTracing Zipkin traecer are 'http' or 'kafka'")
		}
	case "jaeger":
		switch opentracingConfig.Transport {
		case "stdout":
			jaeger, _ := initJaeger(serviceName)
			return jaeger, nil
		default:
			return nil, errors.New("supported transport for OpenTracing Jaeger traecer is 'stdout'")
		}
	default:
		return nil, errors.New("supported implementations for OpenTracing are 'jaeger' or 'zipkin'")
	}
}

// Start implements model.Flow.Start
func (fb *OpenTracingFlow) Start(ctx model.FlowContext) (started bool, taskEntries []*model.TaskEntry) {
	opentracingConfig := readOpentracingContext(ctx)

	if opentracingConfig != nil {
		tracer, err := initTracer(ctx.FlowDefinition().Name(), opentracingConfig)
		if err != nil || tracer == nil {
			log.Warn("Unable to init OpenTracing tracer. Ignoring.")
		} else {
			opentracing.SetGlobalTracer(tracer)

			span := opentracing.StartSpan(ctx.FlowDefinition().Name())
			span.SetTag("type", "flogo:flow")

			// store span in working data to close it later and to pass the span context to activities
			ctx.WorkingData().AddAttr("opentracing-flow-span", data.TypeAny, span)
		}
	} else {
		log.Warn("Unable to init OpenTracing tracer. Ignoring.")
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
