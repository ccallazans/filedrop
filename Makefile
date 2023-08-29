run:
	go run cmd/main.go

lint:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.54.2 golangci-lint run -v
	
migrate-up:
	docker run -v  $(PWD)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

migrate-down:
	docker run -v  $(PWD)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" down 1
