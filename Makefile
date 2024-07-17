## migrate-up
migrate-up: make -migrate -path internal/module/user/migrations -database "postgres://root:root@localhost:5432/currency-api?sslmode=disable" up
rabbit-mq: docker run -d --name rabbitmq --net currency-api-net -p 5672:5672 -p 15672:15672 rabbitmq:3.9.11-management-alpine

## fumpt-all
fumpt-all: gofumpt -w ./..

lint: golintci-linter run