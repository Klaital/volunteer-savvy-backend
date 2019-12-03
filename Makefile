.PHONY: build test clean

build:
	go build -o volunteer-savvy-backend cmd/volunteer-savvy-backend/main.go

test:
	go test github.com/klaital/volunteer-savvy-backend/internal/pkg/organizations
	go test github.com/klaital/volunteer-savvy-backend/internal/pkg/users
	go test github.com/klaital/volunteer-savvy-backend/internal/pkg/sites

clean:
	rm volunteer-savvy-backend
