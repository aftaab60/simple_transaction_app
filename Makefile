docker-db-up:
	cd docker && docker-compose up db_mysql -d
	cd docker && docker-compose up db_pgsql -d

docker-db-down:
	cd docker && docker-compose stop db_mysql  && docker-compose rm -svf  db_mysql
	cd docker && docker-compose stop db_pgsql  && docker-compose rm -svf  db_pgsql

run-migration:
	PGPASSWORD=dev123 psql -U userdev -d db_pgsql -h localhost -p 5432 -f migrations/001_create_base_tables.sql

wire:
	go get github.com/google/wire/cmd/wire && GO111MODULE=on && go run github.com/google/wire/cmd/wire ./internal

# Starts all services and resources required for development
start:
	make docker-db-up
stop:
	make docker-db-down