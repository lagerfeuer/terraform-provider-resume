default: build

.PHONY: build
build:
	go build

.PHONY: install
install:
	go install

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: clean
clean:
	rm -rf ./$(BIN)

