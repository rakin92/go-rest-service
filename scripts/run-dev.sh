#!/bin/sh
app="service"
src="$srcPath/$app/$pkgFile"

printf "\nStart running: $app\n"
# Set all ENV vars for the server to run
export $(grep -v '^#' .env | xargs)
time go run cmd/service/main.go
# This should unset all the ENV vars, just in case.
# unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/''' | xargs)
printf "\nStopped running: $app\n\n"