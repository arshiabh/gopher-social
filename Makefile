
migration:
	@migrate create -ext sql -dir ./cmd/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))
migrate-up:
	@migrate -path ./cmd/migrate/migrations -database postgres://postgres:09300617050@localhost:5432/gopher-database?sslmode=disable up
migrate-down:
	@migrate -path ./cmd/migrate/migrations -database postgres://postgres:09300617050@localhost:5432/gopher-database?sslmode=disable down $(filter-out $@,$(MAKECMDGOALS)) 