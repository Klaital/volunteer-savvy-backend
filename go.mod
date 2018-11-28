module github.com/klaital/volunteer-savvy-backend

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/sirupsen/logrus v1.2.0
)

replace github.com/klaital/volunteer-savvy-backend/internal/pkg/config => ./internal/pkg/config
