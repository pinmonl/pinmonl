NAME := pinmonl
PKG := github.com/pinmonl/pinmonl

.PHONY: fmt
fmt:
	@gofmt -s -l .

.PHONY: image-dev
image-dev:
	@docker build \
		--rm --force-rm \
		-f Dockerfile.dev \
		-t pinmonl/$(NAME):dev .

.PHONY: dshell
dshell:
	@docker run --rm -it \
		-p 8080:8080 \
		-v $(CURDIR):/go/src/$(PKG) \
		--workdir /go/src/$(PKG) \
		pinmonl/$(NAME):dev sh
