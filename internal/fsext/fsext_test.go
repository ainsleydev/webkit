package fsext

import "embed"

//go:embed testdata/*
var testFS embed.FS
