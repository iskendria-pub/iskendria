#!/bin/bash

set -e

rm -rf commandGenerated.go
rm -rf tmp

D=$(pwd)
cd ../generate/command
go build
cd ${D}

../generate/command/command > tmp
gofmt tmp > commandGenerated.go
