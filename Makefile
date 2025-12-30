SHELL := /bin/sh

.PHONY: help
help:
	@printf "%s\n" \
		"Targets:" \
		"  tidy         - go mod tidy for all services" \
		"  test         - go test for all services" \
		"  build        - go build for all services" \
		"  run-catalog  - run catalog service" \
		"  run-order    - run order service"

.PHONY: tidy
tidy:
	@set -e; \
		for d in services/catalog services/order; do \
			( cd $$d && go mod tidy ); \
		done

.PHONY: test
test:
	@set -e; \
		for d in services/catalog services/order; do \
			( cd $$d && go test ./... ); \
		done

.PHONY: build
build:
	@mkdir -p bin
	@set -e; \
		( cd services/catalog && go build -o ../../bin/catalog ./cmd/catalog ); \
		( cd services/order && go build -o ../../bin/order ./cmd/order )

.PHONY: run-catalog
run-catalog:
	@cd services/catalog && PORT=8081 VERSION=dev go run ./cmd/catalog

.PHONY: run-order
run-order:
	@cd services/order && PORT=8082 VERSION=dev go run ./cmd/order
