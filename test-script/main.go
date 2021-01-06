/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
