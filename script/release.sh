#!/bin/sh

cd ./cmd/strikes

#Please set GITHUB_TOKEN in advance. 

set -e
latest_tag=$(git describe --abbrev=0 --tags)
goxc
ghr -u TsuyoshiUshio -r strikes -prerelease $latest_tag  dest/snapshot/