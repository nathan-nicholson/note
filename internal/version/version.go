package version

// Version is the current version of the application
// This will be updated automatically by release-please
var Version = "1.2.0"

// GetVersion returns the current version
func GetVersion() string {
	return Version
}
