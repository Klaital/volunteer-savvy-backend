package auth

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type PasswordCheckResponse struct {
	PasswordMatches bool `db:"password_matches"`
	OrganizationId int64 `db:"organization_id"`
}

func SigninHandler(request *restful.Request, response *restful.Response) {
	var creds Credentials
	decoder := json.NewDecoder(request.Request.Body)
	if decoder == nil {
		log.Error("Failed to construct json decoder")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := decoder.Decode(&creds)
	if err != nil {
		log.Errorf("Failed to decode request body: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Hash the provided password and check for validity
	svcConfig, err := config.GetServiceConfig()
	if err != nil {
		log.Errorf("Service Config not loaded: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	hash, err := creds.HashPassword()
	if err != nil {
		log.Errorf("Failed to hash password: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	var passwordHashMatches PasswordCheckResponse
	row := svcConfig.DatabaseConnection.QueryRow(`SELECT password_digest=? AS password_matches, organization_id FROM users WHERE email=?`, string(hash), creds.Username)
	err = row.Scan(&passwordHashMatches)
	if err != nil {
		log.Errorf("Failed to scan password check results: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if passwordHashMatches.PasswordMatches {
		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username: creds.Username,
			Organization: passwordHashMatches.OrganizationId,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Errorf("Failed to create JWT string: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.AddHeader("token", tokenString)
		err = response.WriteEntity(token)
		if err != nil {
			log.Errorf("Failed to write token as body: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
}