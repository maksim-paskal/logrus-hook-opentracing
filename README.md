## installation
```
go get github.com/maksim-paskal/logrus-hook-opentracing
```

## environment
```
export JAEGER_AGENT_HOST=localhost
export JAEGER_AGENT_PORT=6831
```
## usage

```
package main

import (
	"errors"

	logrushookopentracing "github.com/maksim-paskal/logrus-hook-opentracing"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

var ErrTest error = errors.New("test error")

func main() {
	hook, err := logrushookopentracing.NewHook(logrushookopentracing.Options{})
	if err != nil {
		log.WithError(err).Fatal()
	}

	log.AddHook(hook)

	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.WithError(err).Panic("Could not parse Jaeger env vars")
	}

	cfg.ServiceName = "test-app"
	cfg.Sampler.Type = jaeger.SamplerTypeConst
	cfg.Sampler.Param = 1
	cfg.Reporter.LogSpans = true

	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.WithError(err).Panic("Could not create tracer")
	}
	defer closer.Close()

	span := tracer.StartSpan("main")
	defer span.Finish()

	log.Info("test info")
	log.Warn("test warn")
	log.WithField(logrushookopentracing.SpanKey, span).WithError(ErrTest).Error("test error")
}
```