.PHONY: docker test all lint

all: lint test docker

docker:
	docker build -t cpro29a/mania .
	docker push cpro29a/mania

test:
	GOOGLE_APPLICATION_CREDENTIALS=$(PWD)/credentials.json go test -v -race ./...

test-ci:
	go test -v -race ./...

lint:
	golangci-lint run ./...

checks-ci: lint test-ci
