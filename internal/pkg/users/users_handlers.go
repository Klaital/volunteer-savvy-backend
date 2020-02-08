package users

import (
	"github.com/emicklei/go-restful"
	"github.com/jmoiron/sqlx"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ListUsersRequest struct {
	// Input
	Db *sqlx.DB

	// Output
	Users []User
}

func (request *ListUsersRequest) ListUsers() error {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "ListUsersRequest#ListUsers",
	})

	sqlStmt := request.Db.Rebind(`SELECT * FROM users`)
	users := make([]User, 0)
	err := request.Db.Select(&users, sqlStmt)
	if err != nil {
		logger.WithError(err).Error("Failed to query for users")
		return err
	}

	// Success!
	request.Users = users
	return nil
}
func ListUsersHandler(request *restful.Request, response *restful.Response) {
	logger := logrus.WithFields(logrus.Fields{
		"operation": "ListUsersHandler",
	})

	appConfig, err := config.GetServiceConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load service config")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	requestConfig := ListUsersRequest{
		Db:            appConfig.DatabaseConnection,
	}

	err = requestConfig.ListUsers()
	if err != nil {
		logger.WithError(err).Error("Failed to list users")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = response.WriteEntity(requestConfig.Users)
	if err != nil {
		logger.WithError(err).Error("Failed to serialize users: %+v", requestConfig.Users)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
