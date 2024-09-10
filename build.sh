#!/bin/bash 
set -x
mkdir -p bin
go build -o bin/splitter ./cmd/
