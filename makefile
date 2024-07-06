# Makefile for recursively formatting Go files with goimports

# Define the directories where Go files are located
SRC_DIRS := $(shell find . -type d -not -path '*/\.*')

.PHONY: all fmt

all: fmt

# Format Go files with goimports
fmt:
	@echo "Formatting Go files with goimports..."
	@$(foreach dir,$(SRC_DIRS),goimports -w $(wildcard $(dir)/*.go);)
	@echo "Go files formatted."

