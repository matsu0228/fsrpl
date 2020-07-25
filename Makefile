BINARY_NAME=fsrpl

.PHONY: build lint test test-local emulator realease

build: clean
	go build  -o ./bin/$(BINARY_NAME)

lint:
	golangci-lint run ./...
	go vet ./...

test-local:
	go test ./...

test: lint test-local
	export FIRESTORE_EMULATOR_HOST=0.0.0.0:8080 
	go test ./... -tags=integration

emulator:
	cd emulator && docker-compose up -d && cd ..
	FIRESTORE_EMULATOR_HOST=0.0.0.0:8080 ./emulator/wait.sh

release:
	goreleaser --rm-dist

