.PHONY connect-db
connect-db:
	docker exec -it greenlight_db psql -U test -d greenlight
