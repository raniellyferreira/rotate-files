#!/usr/bin/env bash

set -euo pipefail

VERSION=$(go run src/rotate.go --version | cut -d " " -f 2 | tr -d '[:space:]')

echo "Building binaries"
make build-cross
make dist checksum VERSION="${VERSION}"

echo "Pushing binaries to repository"

aws s3 cp _dist/$VERSION/ s3://awapi-rotate/$VERSION \
    --recursive \
    --exclude "*" \
    --include "rotate-*" \
    --cache-control max-age=3600
