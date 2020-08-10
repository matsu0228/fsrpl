BINARY_NAME=fsrpl

.PHONY: build lint test test-local emulator realease

build: clean
	go build  -o ./bin/$(BINARY_NAME)

lint:
	golangci-lint run ./...

test-local:
	go test ./...

test: lint test-local
	FIRESTORE_EMULATOR_HOST=0.0.0.0:8080 go test ./... -tags=integration

emulator:
	cd emulator && docker-compose up -d && cd ..
	FIRESTORE_EMULATOR_HOST=0.0.0.0:8080 ./emulator/wait.sh

emulator-down:
	cd emulator && docker-compose down && cd ..
	unset FIRESTORE_EMULATOR_HOST


release:
	goreleaser --rm-dist

