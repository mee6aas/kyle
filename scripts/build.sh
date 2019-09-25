#!/usr/bin/env bash

workspace=$(dirname $0)/..
imageRef=mee6aas/runtime-nodejs:latest

cd $workspace

cd ./cmd/kyle-agent

go install -tags netgo -ldflags '-extldflags "-static"'
