#!/bin/bash

set -e

# List of tools to go get
# In the format of "<cli>:<package>" or ":<package>" if no cli
tools=(
  "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
)

tmp_dir=$(mktemp -d -t gotooling-XXX)
echo "Using temp directory ${tmp_dir}"
cp -R boilerplate/golang_support_tools/* $tmp_dir
pushd "$tmp_dir"

for tool in "${tools[@]}"
do
    echo "Installing ${tool}"
    GO111MODULE=on go install $tool
done

popd
