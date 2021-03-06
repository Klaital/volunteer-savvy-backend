package filters

import (
	"context"
	"github.com/emicklei/go-restful"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"os"
)

// SetContextFilter generates a RequestID for this request, and configures a
// logger with it in a field. Both are saved to a context.Context to be passed
// around, and is cached on the Request object as an Attribute.
func SetContextFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// Generate the context
	ctx := GetRequestContext(req)

	// Propagate the request ID to the response as a header
	requestIDptr := ctx.Value("request-id")
	if requestID, ok := requestIDptr.(string); ok && len(requestID) > 0 {
		resp.AddHeader("request-id", requestID)
	}

	// Success!
	chain.ProcessFilter(req, resp)
}

func GetRequestContext(req *restful.Request) context.Context {
	if req == nil {
		return context.Background()
	}

	// Check if the request already has a context
	reqCtx := req.Attribute("ctx")
	if ctx, ok := reqCtx.(context.Context); ok {
		return ctx
	}

	// Initialize the context with a logger preconfigured with a new RequestID
	requestID := req.HeaderParameter("request-id")
	if len(requestID) == 0 {
		requestID = uuid.NewV4().String()
	}

	ctx := context.WithValue(req.Request.Context(), "request-id", requestID)
	logger := GetContextLogger(ctx)
	// cache the preconfigured logger on the context
	ctx = context.WithValue(ctx, "logger", logger)

	return ctx
}

func GetContextLogger(ctx context.Context) *log.Entry {
	// check for a preconfigured logger
	logPtr := ctx.Value("logger")
	if logPtr != nil {
		if logger, ok := logPtr.(*log.Entry); ok {
			return logger
		}
	}

	baseLogger := log.New()

	// Configure the logger from env vars
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel, err := log.ParseLevel(logLevelStr)
	if err != nil {
		//baseLogger.WithError(err).WithField("envLevel", logLevelStr).Warn("Failed to parse log level")
		logLevel = log.DebugLevel
	}
	baseLogger.SetLevel(logLevel)

	logFormat := os.Getenv("LOG_FORMAT")
	switch logFormat {
	case "text":
		baseLogger.SetFormatter(&log.TextFormatter{})
	case "json":
		baseLogger.SetFormatter(&log.JSONFormatter{})
	case "prettyjson":
		baseLogger.SetFormatter(&log.JSONFormatter{
			PrettyPrint: true,
		})
	default:
		baseLogger.SetFormatter(&log.TextFormatter{})
	}

	logger := log.NewEntry(baseLogger)

	// initialize one if needed with standard fields from the context
	if requestID, ok := ctx.Value("request-id").(string); ok && len(requestID) > 0 {
		logger = logger.WithField("RequestID", requestID)
	}

	return logger
}
