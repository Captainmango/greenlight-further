# Greenlight

## Details
Code + read-along with Let's Go Further by Alex Edwards. Some notes are kept in for further reference.

## Usage
There's a Taskfile, so make sure Task is installed. The deps for the project use Docker (the book asks to install things locally, but I said no thank you.)

### Task
You can use `sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d` to install task to a local bin folder if you do not want to install it globally.

Run `docker compose up -d` to start the database

To connect to the database run `task connect-db:greenlight_db:greenlight`. If you named the container and database differently, substitute `greenlight_db` and `greenlight` with the container name and database name respectively.

Install this locally https://github.com/golang-migrate so that you can run migrations. There is a Task command to do so easily `task migrate-db -- up` to bring the database to the latest version.

If we need to rollback for whatever reason, use the following command
```
task migrate-db -- down 1
```

## Tests
This project uses Hurl for e2e API0 contract tests. Install hurl then use `hurl requests/tests/*.hurl --test` to run all tests for the repo.

> Make sure the project is running locally to do this.