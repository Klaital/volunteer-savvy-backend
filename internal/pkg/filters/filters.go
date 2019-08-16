package filters

import (
	"bytes"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
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
		"operation":     "JsonCommonLogger",
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
