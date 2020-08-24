#!/usr/bin/env bash


GO111MODULE=off go get github.com/ka2n/waitport/cmd/waitport/...
export PATH=~/go/bin:$PATH

echo -e "\t$(echo $GOPATH)"
echo -e "\t$(which waitport)"
echo -e "\t$(echo $FIRESTORE_EMULATOR_HOST )"

waitport -listen $FIRESTORE_EMULATOR_HOST -timeout 2m

echo "finish"