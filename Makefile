repo ?= lucky/lucky_box/cms

.PHONY:build-cms
build-cms:
	@cd cms; \
	version=$$(cat version); \
	docker build -t $(repo)_server:$$version .; \
	docker tag $(repo)_server:$$version $(repo)_server:latest


.PHONY: build-frontend
build-frontend:
	@cd view; \
	version=$$(node -p "require('./package.json').version"); \
	docker build -t $(repo)_frontend:$$version .; \
	docker tag $(repo)_frontend:$$version $(repo)_frontend:latest

.PHONY:build
build:
	make build-cms;
	make build-frontend;