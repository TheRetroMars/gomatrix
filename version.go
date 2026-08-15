package main

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var rawVersion string

// Version holds the application version string
// parsed from the embedded VERSION file.
var Version = strings.TrimSpace(rawVersion)
