IMAGE := pinmonl/pinmonl

build-dev:
	@docker build \
		-t $(IMAGE)-dev \
		-f ./Dockerfile.dev \
		.

start-dev:
	@docker run --rm -it \
		--network host \
		-w "$(PWD)" \
		-v "$(PWD):$(PWD)" \
		-v "$(PWD)/.data/go/pkg:/go/pkg" \
		$(IMAGE)-dev sh

fmt:
	@goimports -w $(shell find . -type f -name "*.go" -not -path "./.data/*")

test:
	@go test ./... -v
