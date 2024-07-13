//go:build tools

// Package tools manage tool dependencies via go.mod.

package tools

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser/v2"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
