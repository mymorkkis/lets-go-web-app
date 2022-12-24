package ui

import "embed"

// below is comment directive that instructs Go to embed dirs listed

//go:embed "html" "static"
var Files embed.FS
