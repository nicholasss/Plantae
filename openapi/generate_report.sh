#!/usr/bin/env bash

printf "generating openapi report...\n"
go run github.com/daveshanley/vacuum@latest html-report openapi.yaml
