#!/usr/bin/env bash

set -eo pipefail
DEBUG=${DEBUG:-false}
[[ ${DEBUG} = true ]] && set -x

docker run -it --rm -p 27017:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=mongoadmin \
    -e MONGO_INITDB_ROOT_PASSWORD=secret \
  mongo:latest
