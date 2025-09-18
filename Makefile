repo ?= lucky/lucky_box/cms

.PHONY:build-cms
build-cms:
	@cd cms; \
	version=$$(cat version); \
	docker build -t $(repo):$$version .; \
	docker tag $(repo)/server:$$version $(repo)/server:latest


.PHONY: build-frontend
build-frontend:
	@cd view; \
	version=$$(jq -r '.version' package.json); \
	docker build -t $(repo)_frontend:$$version .; \
	docker tag $(repo)_frontend:$$version $(repo)_frontend:latest
