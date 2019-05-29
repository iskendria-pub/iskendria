#!/bin/bash

set -e

echo "Removing old files..."
rm -rf tmp
rm -rf daoUpdateModificationTimeGenerated.go

D=$(pwd)
echo "Building executable that generates code: modificationTime..."
cd ../generate/dao/modificationTime
go build
cd ${D}
echo "Generating code with executable: modificationTime..."
../generate/dao/modificationTime/modificationTime > tmp
gofmt tmp > daoUpdateModificationTimeGenerated.go
