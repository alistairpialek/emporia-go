go = \
	docker run --rm \
		-v ${PWD}:/go/src \
		-w /go/src \
		${BUILDKITE_PLUGIN_DOCKER_IMAGE} $1

.PHONY: test
test: fmt vet staticcheck
	@echo "--- :go: Test"
	@$(call go,go test -v -cover ./emporia)

.PHONY: fmt
fmt:
	@echo "--- :go: Fmt"
	@$(call go,gofmt -d .)

.PHONY: vet
vet:
	@echo "--- :go: Vet"
	@$(call go,go vet ./...)

.PHONY: staticcheck
staticcheck:
	@echo "--- :go: Staticcheck"
	@$(call go,staticcheck ./...)
