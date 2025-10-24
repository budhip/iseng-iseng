#!/bin/bash

golangci-lint run > lint-report.log

if [ $? -ne 0 ]; then
  go test -cover -race ./tests-scoring/... -v | go-junit-report -set-exit-code > xunitreport.xml
  exit 1
else
  go test -cover -race ./tests-scoring/... -v | go-junit-report -set-exit-code > xunitreport.xml
fi
