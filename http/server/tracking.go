package server

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type KrakenDContextTrackingTypeKey string

const (
	// KrakenDContextOTELStrKey is a special key to be used when there
	// is no way to obtain the span context from an inner context
	// (like when gin has not the fallback option enabled in the engine).
	krakenDContextTrackingStrKey KrakenDContextTrackingTypeKey = "KrakendD-Context-OTEL"
)

type tracking struct {
	startTime time.Time
	ctx       context.Context
	span      trace.Span

	latencyInSecs   float64
	responseSize    int
	responseStatus  int
	responseHeaders map[string][]string
	writeErrs       []error
	endpointPattern string
	isHijacked      bool
	hijackedErr     error
}

func (t *tracking) EndpointPattern() string {
	if len(t.endpointPattern) == 0 {
		if t.isHijacked {
			return "Upgraded Connection"
		}
		return "404 Not Found"
	}
	return t.endpointPattern
}

func newTracking() *tracking {
	return &tracking{
		responseStatus: 200,
	}
}

func fromContext(ctx context.Context) *tracking {
	v := ctx.Value(krakenDContextTrackingStrKey)
	if v != nil {
		t, _ := v.(*tracking)
		return t
	}
	return nil
}

// SetEndpointPattern allows to set the endpoint attribute once it
// has been matched down the http handling pipeline.
func SetEndpointPattern(ctx context.Context, endpointPattern string) {
	if t := fromContext(ctx); t != nil {
		t.endpointPattern = endpointPattern
	}
}

func (t *tracking) Start() {
	t.startTime = time.Now()
}

func (t *tracking) Finish() {
	t.latencyInSecs = float64(time.Since(t.startTime)) / float64(time.Second)
}
