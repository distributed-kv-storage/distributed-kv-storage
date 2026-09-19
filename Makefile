BUF ?= buf
BASE_BRANCH ?= main

.PHONY: ci format lint generate check-generated breaking vet test

ci: lint check-generated vet test

format:
	$(BUF) format -w

lint:
	$(BUF) format --diff --exit-code
	$(BUF) lint

generate:
	$(BUF) generate

check-generated: generate
	git diff --exit-code -- api/gen/go
	test -z "$$(git ls-files --others --exclude-standard api/gen/go)"

breaking:
	$(BUF) breaking --against '.git#branch=$(BASE_BRANCH)'

vet:
	go vet ./...

test:
	go test ./...
