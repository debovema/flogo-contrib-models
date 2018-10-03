package opentracing

import (
	"os"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/debovema/flogo-contrib-models/opentracing/behaviors"
	"github.com/debovema/flogo-contrib-models/opentracing/utils"
)

// log is the default package logger
var log = logger.GetLogger("flowmodel-opentracing")

const (
	MODEL_NAME = "opentracing-model"

	ENV_VARS_PREFIX        = "FLOGO_OPENTRACING_"
	ENV_VAR_IMPLEMENTATION = ENV_VARS_PREFIX + "IMPLEMENTATION"
	ENV_VAR_TRANSPORT      = ENV_VARS_PREFIX + "TRANSPORT"
	ENV_VAR_ENDPOINTS      = ENV_VARS_PREFIX + "ENDPOINTS"
	ENV_VAR_SINGLE_TRACER  = ENV_VARS_PREFIX + "SINGLE_TRACER"
)

func initFromEnvVars() {
	globalOpenTracingImplementation, exists := os.LookupEnv(ENV_VAR_IMPLEMENTATION)
	if !exists {
		return
	}

	log.Infof("Flogo OpenTracing implementation detected: %s.", globalOpenTracingImplementation)

	globalOpenTracingTransport, exists := os.LookupEnv(ENV_VAR_TRANSPORT)
	if !exists {
		log.Errorf("Environment variable %s must be set to initialize OpenTracing tracer.", ENV_VAR_TRANSPORT)
		return
	}
	globalOpenTracingEndpoints, exists := os.LookupEnv(ENV_VAR_ENDPOINTS)
	if !exists {
		log.Errorf("Environment variable %s must be set to initialize OpenTracing tracer.", ENV_VAR_ENDPOINTS)
		return
	}

	behaviors.GlobalConfig = &utils.OpenTracingConfig{Implementation: globalOpenTracingImplementation, Transport: globalOpenTracingTransport, Endpoints: strings.Split(globalOpenTracingEndpoints, ",")}

	globalOpenTracingSingleTracer, exists := os.LookupEnv(ENV_VAR_SINGLE_TRACER)
	if exists {
		useSingleTracer, _ := strconv.ParseBool(globalOpenTracingSingleTracer)
		if useSingleTracer {
			behaviors.GlobalTracer, _ = utils.InitTracer("flogo", behaviors.GlobalConfig)
		}
	}
}

func init() {
	initFromEnvVars()

	model.Register(New())
}

func New() *model.FlowModel {
	m := model.New(MODEL_NAME)
	m.RegisterFlowBehavior(&behaviors.OpenTracingFlow{})
	m.RegisterDefaultTaskBehavior("basic", &behaviors.OpenTracingTask{})
	m.RegisterTaskBehavior("iterator", &behaviors.OpenTracingIteratorTask{})

	return m
}
