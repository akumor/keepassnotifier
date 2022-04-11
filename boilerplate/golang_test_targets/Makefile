
.PHONY: generate
generate: #generate go code
	@boilerplate/golang_test_targets/go-gen.sh

.PHONY: lint
lint: #lints the package for common code smells
	GL_DEBUG=linters_output,env golangci-lint run --deadline=5m --exclude deprecated -v

# If code is failing goimports linter, this will fix.
# skips 'vendor'
.PHONY: goimports
goimports:
	goimports -w $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./pkg/client/*" -not -path "./boilerplate/*")

.PHONY: mod_download
mod_download: #download dependencies (including test deps) for the package
	go mod download

.PHONY: install
install: mod_download

.PHONY: show
show:
	go list -m all

.PHONY: test_unit
test_unit:
	go test -cover ./... -race

.PHONY: test_benchmark
test_benchmark:
	go test -bench . ./...

.PHONY: test_unit_cover
test_unit_cover:
	go test ./... -coverprofile /tmp/cover.out -covermode=count
	go tool cover -func /tmp/cover.out

.PHONY: test_unit_visual
test_unit_visual:
	go test ./... -coverprofile /tmp/cover.out -covermode=count
	go tool cover -html=/tmp/cover.out
	
.PHONY: test_unit_codecov
test_unit_codecov:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic
