build:
	go build -o bin/gts main.go

docker:
	docker build -f Dockerfile \
		-t gql-test-suite:latest \