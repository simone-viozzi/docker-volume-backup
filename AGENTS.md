# AGENT instructions for docker-volume-backup repository

This repository contains Go source code for the main application and a Jekyll
site for documentation.

## Repository layout
- `cmd/` and `internal/` contain the Go source code.
- `docs/` holds the documentation site built using Jekyll.
- `test/` contains integration tests.

## Development guidelines
- Format Go code using `gofmt` (`go fmt ./...`).
- Lint Go code with `golangci-lint run` (configuration in `.golangci.yml`).
- Documentation changes live in the `docs/` directory. Run `bundle install`
  followed by `bundle exec jekyll serve` to preview docs locally.

## Testing
Tests are executed automatically when you open a pull request. You do not
need to run them locally.
