#!/usr/bin/env bash

# Copyright The Rotate Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

VERSION=$(go run ./cmd/rotate version | cut -d " " -f 2 | tr -d '[:space:]')

echo "Building binaries"
make build-cross
make dist checksum VERSION="${VERSION}"

echo "Pushing binaries to repository"

aws s3 cp _dist/$VERSION/ s3://awapi-rotate/$VERSION \
    --recursive \
    --exclude "*" \
    --include "rotate-*" \
    --cache-control max-age=3600

aws s3 cp _dist/$VERSION/ s3://awapi-rotate/latest \
    --recursive \
    --exclude "*" \
    --include "rotate-*" \
    --cache-control max-age=3600
