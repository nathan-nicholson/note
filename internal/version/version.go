package version

// Version is the current version of the application
// This will be injected at build time via ldflags by GoReleaser
// The default value "dev" is used for local development builds
var Version = "dev"

// GetVersion returns the current version
func GetVersion() string {
	return Version
}
