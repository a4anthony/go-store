#include env file and get DB_URL
include .env
export $(shell sed 's/=.*//' .env)

goose-up:
	cd sql/schema &&  goose postgres $(DB_URL) up && cd ../../ && sqlc generate

goose-down:
	cd sql/schema &&  goose postgres $(DB_URL) down

goose-redo:
	cd sql/schema &&  goose postgres $(DB_URL) redo && cd ../../ && sqlc generate

goose-reset:
	cd sql/schema &&  goose postgres $(DB_URL) redo && cd ../../ && sqlc generate


goose-refresh:
	make goose-reset && make goose-up


db-seed:
	go run seeds.go