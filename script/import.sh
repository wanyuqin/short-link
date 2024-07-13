#!/bin/bash

# process_go_files.sh

# Find all Go files and use xargs to process in batches
find . -name '*.go' -print0 | xargs -0 -n 100 go fmt

echo "Go files formatted."
