# create_migration:
	# goose create add_some_table sql

migrate_up:
	GOOSE_DBSTRING="postgres://postgres:password@localhost/outagealertio?sslmode=disable" GOOSE_DRIVER=postgres goose -dir=./migrations up

migrate_down:
	GOOSE_DBSTRING="postgres://postgres:password@localhost/outagealertio?sslmode=disable" GOOSE_DRIVER=postgres goose -dir=./migrations down

run:
	go build -o outagealert && ./outagealert

sqlcgen:
	cd sqlc && sqlc generate
