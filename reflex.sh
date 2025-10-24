#!/bin/bash

# Exit on error
set -e

go install github.com/cespare/reflex@latest

mkdir -p "$HOME/go/bin"

if ! grep -qxF 'export PATH=$PATH:$HOME/go/bin' "$HOME/.bashrc"; then
    echo 'export PATH=$PATH:$HOME/go/bin' >> "$HOME/.bashrc"
    echo "Added to .bashrc"
else
    echo "Already present in .bashrc"
fi

source "$HOME/.bashrc"

nohup reflex -r '\.go$' -s -- sh -c 'golangci-lint run > lint-report.log' > reflex.log 2>&1 &
