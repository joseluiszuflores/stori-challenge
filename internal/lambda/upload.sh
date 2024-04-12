#!/usr/bin/env bash
rm -r bootstrap
rm -r myFunction.zip
GOOS=linux GOARCH=arm64 go build -tags stori  -o bootstrap main.go
zip myFunction.zip bootstrap