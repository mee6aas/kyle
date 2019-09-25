#!/usr/bin/env bash

for var in "$@"
do
    printf '%s\n' "${!var}"
done
