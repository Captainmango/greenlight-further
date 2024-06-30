
.PHONY: connect-db
connect-db:
	docker exec -it greenlight_db psql -U test -d greenlight

.PHONY: migrate-db
migrate-db:
	migrate -path=./migrations -database=postgres://test:password@localhost:5432/greenlight?sslmode=disable up
