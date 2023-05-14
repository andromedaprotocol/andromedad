#!/usr/bin/env bash

set -eo pipefail

mkdir -p ./docs/client
proto_dirs=$(find ./proto -type d -exec sh -c 'find "$0" -maxdepth 1 \( -name "query.proto" -o -name "service.proto" \) | head -n 1' {} \;)
for query_file in $proto_dirs; do
  buf protoc \
    -I "proto" \
    -I "third_party/proto" \
    "$query_file" \
    --swagger_out=./docs/client \
    --swagger_opt=logtostderr=true \
    --swagger_opt=fqn_for_swagger_name=true \
    --swagger_opt=simple_operation_ids=true
done
