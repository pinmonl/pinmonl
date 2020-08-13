define run_app
	$(eval app := $(1))
	$(eval args := $(2))
	@go run ./cmd/$(app)/ $(args)
endef

define build_app
	$(eval app := $(1))
	@pkger -o ./cmd/$(app)
	@go build -o releases/$(app) ./cmd/$(app)
	@rm ./cmd/$(app)/pkged.go
endef

define build_app_image
	$(eval app := $(1))
	@docker build -t pinmonl/$(app) -f docker/Dockerfile.$(app) .
endef

fmt:
	@goimports -w .

run-exchange: export PINMONL_JWT_SECRET = secret
run-exchange: export PINMONL_GIT_DEV = true
run-exchange: export PINMONL_GITHUB_TOKENS = $(shell cat ./.github_token)
run-exchange: export PINMONL_YOUTUBE_TOKENS = $(shell cat ./.youtube_token)
run-exchange: export PINMONL_VERBOSE = 3
run-exchange: args ?= server
run-exchange:
	$(call run_app,exchange,$(args))

run-client: export PINMONL_JWT_SECRET = secret
run-client: export PINMONL_EXCHANGE_ADDRESS = http://localhost:8080
run-client: export PINMONL_EXCHANGE_ENABLED = true
run-client: export PINMONL_VERBOSE = 3
run-client: export PINMONL_WEB_DEVSERVER = http://node:8080
run-client: args ?= server
run-client:
	$(call run_app,pinmonl,$(args))

build-exchange:
	$(call build_app,exchange)

build-client:
	$(call build_app,pinmonl)

build-exchange-image:
	$(call build_app_image,exchange)

build-client-image:
	$(call build_app_image,pinmonl)

