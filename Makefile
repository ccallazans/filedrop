run:
	go run cmd/main.go

migrate:
	migrate -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path ./migrations up

migrate-down:
	migrate -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path ./migrations down