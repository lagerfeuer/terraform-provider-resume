default: build

.PHONY: build
build:
	go build

.PHONY: install
install:
	go install

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean
clean:
	rm -rf ./$(BIN)

# Local API
.PHONY: start-api
start-api:
	docker compose up -d
	docker compose exec app rails db:create db:migrate db:seed

.PHONY: stop-api
stop-api:
	docker compose down

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 RESUME_API_ENDPOINT='http://localhost:3000' RESUME_API_TOKEN='test' \
	go test ./... -v $(TESTARGS) -timeout 120m

