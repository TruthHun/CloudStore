help:
	@echo "test             run test"
	@echo "lint             run lint"
	@echo "example          run examples"

export GO111MODULE=on

.PHONY: test
test:
	go test -race -v -cover -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
	-open cover.html

.PHONY: lint
lint:
	gofmt -s -w .
	goimports -w .
	golint .
	go vet

.PHONY: example
example:
	cd _example && bash test.sh
