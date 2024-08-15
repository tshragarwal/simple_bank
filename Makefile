migrateup:
	migrate -path db/migration -database "postgresql://root:root@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

runtest:
	go test -v -cover ./...

.PHONY: migrateup migratedown sqlc runtest