package tracer

import (
	"fmt"
	"time"

	"go.zenithar.org/pkg/log"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
	"go.uber.org/zap"
)

// Tracer instance
var tracer opentracing.Tracer

// SetTracer can be used by unit tests to provide a NoopTracer instance. Real users should always
// use the InitTracing func.
func SetTracer(initializedTracer opentracing.Tracer) {
	tracer = initializedTracer
}

// Sampling Server URL
var samplingServerURL string

// SetSamplingServerURL defines the global server url
func SetSamplingServerURL(url string) {
	samplingServerURL = url
}

// InitTracing connects the calling service to Zipkin and initializes the tracer.
func InitTracing(serviceName string, logger log.LoggerFactory) opentracing.Tracer {
	logger.Bg().Debug("Initializing tracer", zap.String("svc", serviceName))

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)

	// Zipkin shares span ID between client and server spans; it must be enabled via the following option.
	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	// sender, err := jaeger.NewUDPTransport("jaeger-agent.istio-system:5775", 0)
	sender, err := jaeger.NewUDPTransport(samplingServerURL, 0)
	if err != nil {
		logger.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}

	tracer, _ := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second)),
		injector,
		extractor,
		zipkinSharedRPCSpan,
		jaeger.TracerOptions.Logger(jaegerLoggerAdapter{logger.Bg()}),
	)
	return tracer
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
