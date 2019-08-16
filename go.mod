module github.com/klaital/volunteer-savvy-backend

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/emicklei/go-restful v2.9.6+incompatible
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170926063155-7524189396c6
	github.com/golang-migrate/migrate/v4 v4.5.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.0.0
	github.com/sirupsen/logrus v1.4.1
)

replace github.com/klaital/volunteer-savvy-backend/internal/pkg/config => ./internal/pkg/config

replace github.com/klaital/volunteer-savvy-backend/internal/pkg/users => ./internal/pkg/users
