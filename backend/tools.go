//go:build tools

// Package tools pins tool dependencies so they are recorded in go.mod / go.sum.
package tools

import (
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)
