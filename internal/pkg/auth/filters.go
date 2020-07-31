package auth

import (
	"context"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"net/http"
	"strings"
)

// ValidJwtFilter ensures that an API request is made with a valid, signed bearer token
func ValidJwtFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := filters.GetRequestContext(req)
	logger := filters.GetContextLogger(ctx)
	authHeader := req.HeaderParameter("Authorization")
	authHeaderTokens := strings.Split(authHeader, " ")
	if len(authHeaderTokens) != 2 {
		logger.WithField("Authorization", authHeader).Debug("invalid token")
		resp.WriteHeader(http.StatusForbidden)
		return
	}

	tokenType, jwtRaw := authHeaderTokens[0], authHeaderTokens[1]
	if !strings.EqualFold(tokenType, "Bearer") {
		logger.WithField("TokenType", tokenType).Debug("invalid token type")
		resp.WriteHeader(http.StatusForbidden)
		return
	}

	cfg, err := config.GetServiceConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to get service configuration for JWT filter")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	userGuid, token, err := ParseJwt(jwtRaw, cfg.JwtPublicKey)
	logger = logger.WithField("jwt.sub", userGuid)
	if err != nil {
		logger.WithError(err).Debug("Failed to parse the JWT")
		resp.WriteHeader(http.StatusForbidden)
		return
	}
	// Update the request context with the logged-in user ID
	ctx = context.WithValue(ctx, "logger", logger)
	req.SetAttribute("ctx", ctx)
	req.SetAttribute("jwt", token)
	req.SetAttribute("jwt.sub", userGuid)

	chain.ProcessFilter(req, resp)
}
