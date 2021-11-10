#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

VERSION="${1}"

if [[ -z "${VERSION}" ]]; then
  echo "Usage: push_tag.sh <version>"
  exit 1
fi
PREVIOUS_VERSION=$( git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')

sed -i "" "s/go-logging@v${PREVIOUS_VERSION}/go-logging@v${VERSION}/" README.md
git add README.md
git commit -m "bump version to v${VERSION}"
git push

git tag -a "v${VERSION}" -m "version ${VERSION}"
git push origin "v${VERSION}"
