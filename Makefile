include .env
export

goose_up:
	goose -dir=./pkg/storage/postgres/migrations up

sqlc_generate: goose_up
	sqlc generate -f ./pkg/storage/postgres/sqlc.yaml

dev: sqlc_generate
	go run main.go
