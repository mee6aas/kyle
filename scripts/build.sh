#!/usr/bin/env bash

workspace=$(dirname $0)/..

cd $workspace

cd ./cmd/kyle-agent

CGO_ENABLED=0 go install -tags netgo -ldflags '-extldflags "-static -w"'
