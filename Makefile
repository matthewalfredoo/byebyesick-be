include .env
include .dockerimages

postgres:
	docker run --name ${POSTGRES_CONTAINER_NAME} -p ${DB_PORT}:5432 -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d ${POSTGRES_IMAGE}

rundb:
	docker start postgres16

createdb:
	docker exec -it ${POSTGRES_CONTAINER_NAME} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

dropdb:
	docker exec -it ${POSTGRES_CONTAINER_NAME} dropdb ${DB_NAME}

migrateup:
	migrate -path app/appdb/migration --database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migratedown:
	migrate -path app/appdb/migration --database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
