build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/main.go

run:
	go run cmd/*

migrate-up:
	docker run -v  $(PWD)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

migrate-down:
	docker run -v  $(PWD)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" down 1
