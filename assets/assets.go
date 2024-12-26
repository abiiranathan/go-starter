package assets

import "embed"

// Static contains all static files.
// We are creating a seperate package for assets to embed the static files
// because we can't use relative paths go:embed directive.
// This also helps if we want to access the static files in other packages.
//
//go:embed static/*
var Static embed.FS
