# OpenTracing model for Flogo

This activity adds OpenTracing instrumentation at Flogo engine level and tracers implementations
(Zipkin over HTTP or Kafka, Jaeger).

## Installation

### Flogo Web

This model is not available with the Flogo Web UI

### Flogo CLI

In the directory of a Flogo project (with a *flogo.json* file) :

```bash
flogo install github.com/debovema/flogo-contrib-models/opentracing
```

#### Patch flogo-contrib and flogo-lib

This model requires some little updates in flogo-contrib and flogo-lib which are not yet merged into TIBCOSoftware
repositories.
A script is provided to perform the operation.

In the directory of the Flogo project (with a *flogo.json* file) :

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/debovema/flogo-contrib-models/master/opentracing/patch-vendor.sh)"
```

## Usage

In the *flogo.json*, replace 

```json
  "resources": [
    {
      "id": "flow:sample_flow",
      "data": {
        "name": "SampleFlow",
        "tasks": [
          {
            "id": "log_2",
            "name": "Log Message",
            "description": "Simple Log Activity",
            "activity": {
              "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
              "input": {
                "message": "Simple Log",
                "flowInfo": "false",
                "addToFlow": "false"
              }
            }
          }
        ]
      }
    }
  ]
```

by 

```json
  "resources": [
    {
      "id": "flow:sample_flow",
      "data": {
        "name": "SampleFlow",
        "model": "github.com/square-it/flogo-contrib-models/opentracing",
        "attributes": [
          {
            "name": "opentracing-config-http",
            "type": "any",
            "value": {
              "implementation": "zipkin",
              "transport": "http",
              "endpoints": [
                "http://127.0.0.1:9411/api/v1/spans"
              ]
            }
          }
        ],
        "tasks": [
          {
            "id": "log_2",
            "name": "Log Message",
            "description": "Simple Log Activity",
            "activity": {
              "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
              "input": {
                "message": "Simple Log",
                "flowInfo": "false",
                "addToFlow": "false"
              }
            }
          }
        ]
      }
    }
  ]
```

Replace *127.0.0.1* by the actual IP of the Zipkin collector.