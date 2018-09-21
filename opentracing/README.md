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


