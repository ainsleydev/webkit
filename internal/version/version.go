package version

// Version is the current version number that Webkit is running on.
// This value is injected by GoReleaser during builds via ldflags.
var Version = "dev"

// BuildInfo contains additional build metadata injected by GoReleaser.
var (
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "unknown"
)

// Info returns a formatted string with all version information.
func Info() string {
	return "WebKit " + Version + "\n" +
		"  Commit:   " + Commit + "\n" +
		"  Built:    " + Date + "\n" +
		"  Built by: " + BuiltBy
}
