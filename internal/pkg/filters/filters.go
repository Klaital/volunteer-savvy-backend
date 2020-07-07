package filters

import (
	"bytes"
	"context"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
)

func getOriginatingIP(req *restful.Request) (proxyExists bool, originatingIP string) {
	// obtain the ip the request originated from which is possibly not
	// the requesting ip inside the load balancer.
	originatingIpSet, proxyExists := req.Request.Header["X-Forwarded-For"]

	if proxyExists {
		originatingIP = strings.Join(originatingIpSet, " ")
	} else {
		originatingIP = strings.Split(req.Request.RemoteAddr, ":")[0]
	}

	return
}

func JSONCommonLogger(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)

	_, originatingIP := getOriginatingIP(req)

	var requestBody = ""
	if req.Request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(req.Request.Body)
		req.Request.Body.Close()
		req.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		requestBody = string(bodyBytes)
	}

	fields := log.Fields{
		"operation":     "JSONCommonLogger",
		"originatingIP": originatingIP,
		"time":          time.Now().Format("02/Jan/2006:15:04:05.000 -0700"),
		"method":        req.Request.Method,
		"requestURI":    req.Request.URL.RequestURI(),
		"proto":         req.Request.Proto,
		"statusCode":    resp.StatusCode(),
		"contentLength": resp.ContentLength(),
		"requestBody":   requestBody,
	}

	if len(requestBody) == 0 {
		delete(fields, "requestBody")
	}

	logger := log.WithFields(fields)

	// workaround for log spam: health check requests are debug-level
	// log lines
	cfg, _ := config.GetServiceConfig()
	if strings.HasPrefix(req.Request.URL.RequestURI(), cfg.HealthCheckPath) {
		logger.Debugln()
	} else {
		logger.Infoln()
	}
}

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

	ctx := context.WithValue(context.Background(), "request-id", requestID)
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

	logger := log.NewEntry(log.New())

	// initialize one if needed with standard fields from the context
	if requestID, ok := ctx.Value("request-id").(string); ok && len(requestID) > 0 {
		logger = logger.WithField("RequestID", requestID)
	}

	return logger
}