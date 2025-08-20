# Those are callable targets
TASKS = $(shell cd build && go run . --list)

.PHONY: all
all: build

.PHONY: $(TASKS)
$(TASKS):
	@cd build && go run . $(ARGS) $@

.PHONY: help
help:
	@cd build && go run . --help