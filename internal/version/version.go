package version

// Version is the current version number that Webkit is running on.
// This value is generated at build time via generate.go.
var Version = getVersion()

// BuildInfo contains additional build metadata.
var (
	Commit  = getCommit()
	Date    = getDate()
	BuiltBy = getBuiltBy()
)

// Info returns a formatted string with all version information.
func Info() string {
	return "WebKit " + Version + "\n" +
		"  Commit:   " + Commit + "\n" +
		"  Built:    " + Date + "\n" +
		"  Built by: " + BuiltBy
}

// getVersion returns the version from generated constants or a default.
func getVersion() string {
	if v := generatedVersion; v != "" {
		return v
	}
	return "v0.0.1-dev"
}

// getCommit returns the commit from generated constants or a default.
func getCommit() string {
	if c := generatedCommit; c != "" {
		return c
	}
	return "none"
}

// getDate returns the date from generated constants or a default.
func getDate() string {
	if d := generatedDate; d != "" {
		return d
	}
	return "unknown"
}

// getBuiltBy returns the built by info from generated constants or a default.
func getBuiltBy() string {
	if b := generatedBuiltBy; b != "" {
		return b
	}
	return "local"
}
