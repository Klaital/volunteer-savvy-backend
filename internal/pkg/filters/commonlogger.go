package filters

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
)

// JsonLoggingFilter generates log lines for requests in JSON format.
func JsonLoggingFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	_, originatingIP := getOriginatingIP(req)

	ctx := GetRequestContext(req)
	logger := GetContextLogger(ctx)

	// W3C Trace Context Traceparent header
	// The traceparent header represents the incoming request in a tracing system in a common format, understood by all vendors.
	traceparent := req.HeaderParameter("traceparent")

	// We do this before calling ProcessFilter because ProcessFilter consumes the request body.
	var requestBody = ""
	if req.Request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(req.Request.Body)
		requestBody = string(bodyBytes)
		req.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	chain.ProcessFilter(req, resp)

	// TODO: Would be nice to log the response body.

	fields := logrus.Fields{
		"operation":     "NCSACommonLogFormatLogger",
		"originatingIP": originatingIP,
		"time":          time.Now().Format("02/Jan/2006:15:04:05.000 -0700"),
		"method":        req.Request.Method,
		"requestURI":    req.Request.URL.RequestURI(),
		"proto":         req.Request.Proto,
		"statusCode":    resp.StatusCode(),
		"contentLength": resp.ContentLength(),
		"traceparent":   traceparent,
		"requestBody":   requestBody,
	}

	if len(requestBody) == 0 {
		delete(fields, "requestBody")
	}

	logMsgNCSACLF := fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %d",
		originatingIP,
		time.Now().Format("02/Jan/2006:15:04:05.000 -0700"),
		req.Request.Method,
		req.Request.URL.RequestURI(),
		req.Request.Proto,
		resp.StatusCode(),
		resp.ContentLength(),
	)

	logger = logger.WithFields(fields)

	// Health Check requests are debug-level log lines to reduce logspam in non-test realms
	if isHealthRequest(req) {
		logger.Debugln(logMsgNCSACLF)
	} else {
		logger.Infoln(logMsgNCSACLF)
	}
}

func isHealthRequest(req *restful.Request) bool {
	return strings.HasPrefix(req.Request.URL.RequestURI(), "/GetServiceStatus") || strings.HasPrefix(req.Request.URL.RequestURI(), "/healthz")
}

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
