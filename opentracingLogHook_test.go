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
	"errors"
	"testing"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

var ErrTest error = errors.New("test error")

func TestHook(t *testing.T) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("testHook")

	defer span.Finish()

	t.Parallel()

	hook, err := NewHook(Options{})
	if err != nil {
		t.Fatal(err)
	}

	log.AddHook(hook)

	log.Info("test info")
	log.Warn("test warn")
	log.WithField(SpanKey, span).WithError(ErrTest).Error("test error")
}
