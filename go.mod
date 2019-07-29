module github.com/klaital/volunteer-savvy-backend

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/sirupsen/logrus v1.2.0
)

replace github.com/klaital/volunteer-savvy-backend/internal/pkg/config => ./internal/pkg/config

replace github.com/klaital/volunteer-savvy-backend/internal/pkg/users => ./internal/pkg/users
