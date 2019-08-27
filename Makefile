BINARY_NAME=gfs

build: clean
	go build  -o ./bin/$(BINARY_NAME)

test:
		go test ./... -v
		golint ./...
		go vet ./...
		errcheck ./...

clean:
		rm -f ./bin/*

# copy to $GOBIN
install: build
	cp -f ./bin/$(BINARY_NAME) $(GOBIN)/

# build release binary
release: clean
	GOOS=darwin  GOARCH=amd64  go build -o ./bin/$(BINARY_NAME) &&  zip MacOS.zip   ./bin && rm -f ./bin
	GOOS=linux   GOARCH=amd64  go build -o ./bin/$(BINARY_NAME) &&  zip Linux.zip   ./bin && rm -f ./bin
	GOOS=windows GOARCH=amd64  go build -o ./bin/$(BINARY_NAME) &&  zip Windows.zip ./bin && rm -f ./bin
