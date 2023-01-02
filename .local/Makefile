.PHONY: run
run:
	go run .

.PHONY: test
test: fmt vet staticcheck
	go test -v -cover ./emporia

# Like go fmt, but just display diffs.
.PHONY: fmt
fmt:
	gofmt -d .

.PHONY: vet
vet:
	go vet ./...

# https://github.com/golang/vscode-go/wiki/tools#staticcheck is installed by vscode go extension.
.PHONY: staticcheck
staticcheck:
	staticcheck ./...
