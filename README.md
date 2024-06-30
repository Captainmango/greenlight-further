# Greenlight

## Details
Code + read-along with Let's Go Further by Alex Edwards. Some notes are kept in for further reference.

## Usage
There's a Makefile. The deps for the project use Docker (the book asks to install things locally, but I said no thank you.)

Run `docker compose up -d` to start the database

To connect to the database run `make connect-db`

Install this locally https://github.com/golang-migrate so that you can run migrations. There is a Make command to do so easily `make migrate-db`

If we need to rollback for whatever reason, use the following command
```
migrate -path=./migrations -database=postgres://test:password@localhost:5432/greenlight?sslmode=disable down 1
```
