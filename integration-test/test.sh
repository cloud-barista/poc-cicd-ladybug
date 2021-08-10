#!/bin/bash

source ./test.env

rm -f profile.cov
rm -rf meta_db

go test -p 1 -v -coverpkg=$(go list ../... | grep -v integration-test | grep -v protobuf | tr "\n" ",")  -coverprofile=profile.cov ./...

echo "============================================="
totalCoverage=$(go tool cover -func profile.cov | grep total | awk '{print $3}')
echo "Total Coverage : $totalCoverage"
echo "============================================="
