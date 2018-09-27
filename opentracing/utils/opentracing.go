package utils

import (
	"errors"
	"fmt"
	"io"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"

	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

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

func ReadOpentracingContext(ctx model.FlowContext) *OpenTracingConfig {
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

func InitTracer(serviceName string, opentracingConfig *OpenTracingConfig) (opentracing.Tracer, error) {
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
