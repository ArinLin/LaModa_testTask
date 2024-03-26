tests:
	go test -count=1 ./internal/store/...

coverage:
	go test -count=1 -coverprofile=coverage.out ./internal/store/...
	go tool cover -html="coverage.out"

up:
	docker compose --env-file deploy/.env -f deploy/docker-compose.yml -p lamoda_inventory_hub up -d --build

down:
	docker compose -p lamoda_inventory_hub down

logs:
	docker compose -p lamoda_inventory_hub logs

fmt:
	go fmt ./...

swagger-generate:
	cd ./api/hub/rest/handlers && swag init -g ../../../../cmd/main.go --pd -o ../../../doc

wire:
	go run -mod=mod github.com/google/wire/cmd/wire ./api/hub/inject/

