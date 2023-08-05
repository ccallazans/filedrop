run:
	go run cmd/main.go
	
migrate-up:
	docker run -v  $(PWD)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

migrate-down:
	docker run -v  $(PWD)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" down 1
