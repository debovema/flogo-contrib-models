package opentracing

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	simple_behaviors "github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/square-it/flogo-contrib-models/opentracing/behaviors"
)

// log is the default package logger
var log = logger.GetLogger("flowmodel-opentracing")

const (
	MODEL_NAME = "github.com/square-it/flogo-contrib-models/opentracing"
)

func init() {
	model.Register(New())
}

func New() *model.FlowModel {
	m := model.New(MODEL_NAME)
	m.RegisterFlowBehavior(&behaviors.OpenTracingFlow{})
	m.RegisterDefaultTaskBehavior("basic", &behaviors.OpenTracingTask{})
	m.RegisterTaskBehavior("iterator", &simple_behaviors.IteratorTask{})

	return m
}
