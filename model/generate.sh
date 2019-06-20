#!/bin/bash

set -e

rm -rf modelGenerated.go
rm -rf tmp
rm -rf ../generate/model/model

D=$(pwd)
cd ../generate/model
go build
cd ${D}

../generate/model/model > tmp
gofmt tmp > modelGenerated.go
