#!/bin/bash

set -e

rm -rf commandGenerated.go
rm -rf commandUpdateModificationTimeGenerated.go
rm -rf tmp

D=$(pwd)
cd ../generate/command/command
go build
cd ${D}
cd ../generate/command/modificationTime
go build
cd ${D}
cd ../generate/command/unmarshalledState
go build
cd ${D}

../generate/command/command/command > tmp
gofmt tmp > commandGenerated.go
../generate/command/modificationTime/modificationTime > tmp
gofmt tmp > commandUpdateModificationTimeGenerated.go
../generate/command/unmarshalledState/unmarshalledState > tmp
gofmt tmp > unmarshalledStateGenerated.go
