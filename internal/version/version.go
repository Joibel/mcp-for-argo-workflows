// Package version provides build version information.
package version

// These variables are set at build time via ldflags.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)
