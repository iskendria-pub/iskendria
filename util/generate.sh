#!/bin/bash

./clean.sh

set -e

D=$(pwd)
cd ../generate/util
go build
cd ${D}

../generate/util/util > tmp
gofmt tmp > utilGenerated.go

