.PHONY: build test clean

build:
	go build -o volunteer-savvy-backend cmd/volunteer-savvy-backend/main.go

testdb:
	docker run --rm -p "5556:5432" -e "POSTGRES_USER=vstester" -e "POSTGRES_PASSWORD=vstester" -e "POSTGRES_DB=vstest" --name "vstest" timms/postgres-logging:10.3

test:
	go test github.com/klaital/volunteer-savvy-backend/internal/pkg/organizations
	go test github.com/klaital/volunteer-savvy-backend/internal/pkg/users
	go test github.com/klaital/volunteer-savvy-backend/internal/pkg/sites

clean:
	rm volunteer-savvy-backend
