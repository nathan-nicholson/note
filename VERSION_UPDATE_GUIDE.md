# Version and Update System Guide

This guide explains how the version and update system works in the `note` CLI.

## Overview

The `note` CLI includes built-in version management and automatic update checking:

- **Version Command**: Display current version and build information
- **Update Command**: Check for and install updates automatically

## Components

### 1. Version Management (`internal/version/version.go`)

```go
var Version = "0.1.0"
```

- The version is stored in a simple Go file
- GoReleaser automatically updates this during releases using ldflags
- Build command: `-ldflags="-X github.com/nathan-nicholson/note/internal/version.Version={{.Version}}"`

### 2. Update Checker (`internal/update/checker.go`)

Provides utilities for:
- Checking GitHub releases API for latest version
- Comparing semantic versions
- Detecting installation method (Homebrew, go install, binary)
- Providing appropriate update instructions
- Performing automatic updates (where possible)

### 3. Commands (`cmd/version.go`, `cmd/update.go`)

- `note version` - Display current version
- `note update` - Interactive update
- `note update --check` - Check only, don't install
- `note update -y` - Update without confirmation

## How It Works

### Version Display

```bash
$ note version
note version 0.1.0
Built with go1.24.2 for darwin/arm64
```

Shows:
- Current version (injected at build time)
- Go version used to build
- Platform and architecture

### Update Check

```bash
$ note update --check
Checking for updates...

Current version: 0.1.0
Latest version:  0.2.0

⚠ A new version is available!

Release notes:
------------------------------------------------------------
## What's New
- Added project archiving
- Improved performance
- Bug fixes
------------------------------------------------------------

To update, run:
  brew upgrade nathan-nicholson/tap/note
```

The update checker:
1. Queries GitHub API: `https://api.github.com/repos/nathan-nicholson/note/releases/latest`
2. Parses the release JSON
3. Compares versions using semantic versioning
4. Displays release notes
5. Provides update instructions

### Automatic Update

```bash
$ note update
Checking for updates...

Current version: 0.1.0
Latest version:  0.2.0

⚠ A new version is available!

Would you like to update now? [y/N]: y

Updating note...
✓ Successfully updated to version 0.2.0!

Please restart note to use the new version.
```

For Homebrew installations:
- Runs: `brew upgrade nathan-nicholson/tap/note`

For go install:
- Runs: `go install github.com/nathan-nicholson/note@latest`

For binary installations:
- Provides download link

### Installation Method Detection

The update system automatically detects how `note` was installed:

#### Homebrew Detection
- Checks if binary is in Homebrew Cellar path
- Runs `brew --prefix note` to verify
- Enables automatic updates via `brew upgrade`

#### Go Install Detection
- Checks if binary is in `$GOPATH/bin` or `$GOBIN`
- Checks default `~/go/bin` location
- Enables automatic updates via `go install`

#### Binary Install (Fallback)
- Any other installation method
- Provides manual download instructions

## Version Comparison

Versions are compared using semantic versioning (semver):

```
1.2.3
│ │ │
│ │ └─ Patch (bug fixes)
│ └─── Minor (new features, backwards compatible)
└───── Major (breaking changes)
```

Examples:
- `1.2.4 > 1.2.3` (newer patch)
- `1.3.0 > 1.2.9` (newer minor)
- `2.0.0 > 1.9.9` (newer major)

## Integration with Release Process

### Automated Version Injection

When GoReleaser creates a release:

```yaml
builds:
  - ldflags:
      - -X github.com/nathan-nicholson/note/internal/version.Version={{.Version}}
```

This replaces the version at build time with the actual release version.

### Release Flow

```
1. Commit with conventional commits
   └─> release-please creates PR with version bump
       └─> Merge PR creates release (e.g., v0.2.0)
           └─> GoReleaser builds with version 0.2.0 injected
               └─> Users can check for updates
                   └─> Update command sees new version
```

## Testing

The update package includes comprehensive tests:

```bash
# Run update tests
go test ./internal/update/...

# Test version comparison
go test ./internal/update/... -run TestCompareVersions

# Test installation method detection
go test ./internal/update/... -run TestDetectInstallMethod
```

## Usage Examples

### Check current version
```bash
note version
```

### Check for updates without installing
```bash
note update --check
```

### Interactive update
```bash
note update
# Prompts for confirmation
```

### Automatic update (skip prompt)
```bash
note update -y
# Updates immediately if newer version available
```

### Scripted update check
```bash
# Exit code 0 if up to date, 1 if update available
note update --check > /dev/null 2>&1
if [ $? -eq 0 ]; then
  echo "Up to date"
else
  echo "Update available"
fi
```

## Error Handling

### No Internet Connection
```
Error: failed to check for updates: dial tcp: lookup api.github.com: no such host
```

### GitHub API Rate Limiting
```
Error: GitHub API returned status 403
```
Solution: Wait an hour or authenticate with GitHub token

### No Releases Available
```
Error: GitHub API returned status 404
```
This is normal for new projects before the first release

### Update Fails
```
Error: failed to update: exit status 1
```
The command output will show the specific error from Homebrew or go install

## Troubleshooting

### "Already on latest version" but version seems old

Check the actual installed version:
```bash
note version
which note
```

If multiple installations exist:
```bash
# macOS/Linux
ls -la $(which note)

# Check Homebrew specifically
brew list note
```

### Update command not working

1. Check internet connection
2. Verify GitHub releases exist: https://github.com/nathan-nicholson/note/releases
3. Try manual update:
   - Homebrew: `brew upgrade nathan-nicholson/tap/note`
   - Go: `go install github.com/nathan-nicholson/note@latest`

### Version shows "0.1.0" in releases

This means version injection failed. The version should be injected by GoReleaser.

Check `.goreleaser.yml`:
```yaml
ldflags:
  - -X github.com/nathan-nicholson/note/internal/version.Version={{.Version}}
```

## Future Enhancements

Potential improvements:

1. **Automatic update checks**: Check on startup (with cache/throttling)
2. **Update notifications**: Show message if update available
3. **Changelog viewing**: `note changelog` command
4. **Version history**: `note version --all` to list releases
5. **Rollback**: `note update --version 0.1.0` to install specific version
6. **Update preferences**: Config file for auto-update behavior

## Security Considerations

The update system:
- ✅ Uses HTTPS for all API calls
- ✅ Verifies checksums (when using GoReleaser downloads)
- ✅ Relies on GitHub's security for releases
- ✅ Never executes arbitrary code from responses
- ✅ Shows release notes before updating
- ✅ Requires user confirmation (unless -y flag)

For additional security, users can:
- Pin to specific version with Homebrew
- Verify release signatures (future enhancement)
- Build from source

## Additional Resources

- [Semantic Versioning](https://semver.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Releases API](https://docs.github.com/en/rest/releases)
- [Homebrew Formula](https://docs.brew.sh/Formula-Cookbook)
