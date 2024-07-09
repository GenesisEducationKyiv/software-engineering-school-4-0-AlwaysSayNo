## migrate-up
migrate-up: make -migrate -path internal/module/user/migrations -database "postgres://root:root@localhost:5432/currency-api?sslmode=disable" up

## fumpt-all
fumpt-all: gofumpt -w ./..

lint: golintci-linter run