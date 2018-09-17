package opentracing

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/model"
	"testing"
)

func TestRegistered(t *testing.T) {
	act := model.Get("flowmodel-opentracing")

	if act == nil {
		t.Error("Model Not Registered")
		t.Fail()
		return
	}
}
