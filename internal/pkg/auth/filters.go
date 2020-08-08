package auth

import (
	"context"
	"crypto/rsa"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/filters"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	"net/http"
	"strings"
)

func  (authConfig AuthConfig) extractJWT(req *restful.Request) (*Claims) {
	authHeader := req.HeaderParameter("Authorization")
	authHeaderTokens := strings.Split(authHeader, " ")
	if len(authHeaderTokens) != 2 {
		return nil
	}

	tokenType, jwtRaw := authHeaderTokens[0], authHeaderTokens[1]
	if !strings.EqualFold(tokenType, "Bearer") {
		return nil
	}

	token, err := ParseJwt(jwtRaw, authConfig.PublicKey)
	if err != nil {
		return nil
	}

	claims, ok := token.Claims.(Claims)
	if !ok {
		return nil
	}
	return &claims
}

// AuthConfig is used as a reciever for auth-related filters in order to pass
// in state such as Public Keys that are typically pulled from env vars. In
// test they will likely be hardcoded.
type AuthConfig struct {
	PublicKey *rsa.PublicKey
}
// ValidJwtFilter ensures that an API request is made with a valid, signed bearer token
func (authConfig AuthConfig) ValidJwtFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
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
	_, publicKey := cfg.GetJWTKeys()
	token, err := ParseJwt(jwtRaw, publicKey)
	claims, ok := token.Claims.(Claims)
	if !ok {
		logger.Error("Failed to parse claims")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	logger = logger.WithField("jwt.sub", claims.Subject)
	if err != nil {
		logger.WithError(err).Debug("Failed to parse the JWT")
		resp.WriteHeader(http.StatusForbidden)
		return
	}
	// Update the request context with the logged-in user ID
	ctx = context.WithValue(ctx, "logger", logger)
	req.SetAttribute("ctx", ctx)
	req.SetAttribute("jwt", token)
	req.SetAttribute("jwt.sub", claims.Subject)

	chain.ProcessFilter(req, resp)
}

// RequiresSuperAdminFilter ensures that the logged-in user has SuperAdmin permissions.
// You should add ValidJwtFilter before this one in the chain.
func (authConfig AuthConfig) RequiresSuperAdminFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	jwtClaims := authConfig.extractJWT(req)
	for _, superRole := range jwtClaims.Roles[0] {
		if superRole.Role == users.SiteAdmin {
			chain.ProcessFilter(req, resp)
			return
		}
	}
	// If we got here, the filter fails
	resp.WriteHeader(http.StatusForbidden)
}
