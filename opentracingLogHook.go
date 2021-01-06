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
package logrushookopentracing

import (
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type Hook struct {
	logLevels []log.Level
}

type Options struct {
	LogLevels []log.Level
}

const SpanKey = "span"

// create new Hook.
func NewHook(options Options) (*Hook, error) {
	hook := Hook{}

	hook.logLevels = options.LogLevels

	if hook.logLevels == nil {
		hook.logLevels = []log.Level{
			log.ErrorLevel,
			log.FatalLevel,
			log.WarnLevel,
			log.PanicLevel,
		}
	}

	return &hook, nil
}

// func log.Hook.Levels.
func (hook *Hook) Levels() []log.Level {
	return hook.logLevels
}

//nolint:funlen
// func log.Hook.Fire.
func (hook *Hook) Fire(entry *log.Entry) error {
	var err error = nil

	if dataErr, ok := entry.Data[log.ErrorKey].(error); ok && dataErr != nil {
		err = dataErr
	}

	if span, ok := entry.Data[SpanKey].(opentracing.Span); ok && span != nil && entry.Level >= log.ErrorLevel {
		span.SetTag("error", true)

		if err != nil {
			span.LogKV("error", err)
		} else {
			span.LogKV("error", entry.Message)
		}
	}

	return nil
}
