#!/bin/bash
rm -rf ./app
mkdir -p ./app
for CMD in `ls ./cmd`; do
  go build -o ./app/$CMD ./cmd/$CMD
done
