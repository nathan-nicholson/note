# Homebrew Publishing Setup

This guide explains how to set up automated Homebrew publishing for the `note` CLI tool.

## Overview

The release process uses:
- **release-please**: Manages versioning and creates GitHub releases
- **GoReleaser**: Builds cross-platform binaries and publishes to Homebrew
- **Homebrew Tap**: A custom repository for the Homebrew formula

## Setup Steps

### 1. Create a Homebrew Tap Repository

Create a new GitHub repository named `homebrew-tap`:

```bash
# On GitHub, create a new repository: homebrew-tap
# Then clone it locally:
git clone https://github.com/nathan-nicholson/homebrew-tap.git
cd homebrew-tap

# Create the Formula directory
mkdir -p Formula

# Create a basic README
cat > README.md << 'EOF'
# Homebrew Tap for nathan-nicholson

This is a custom Homebrew tap for installing tools developed by nathan-nicholson.

## Installation

```bash
brew tap nathan-nicholson/tap
brew install note
```

## Available Formulae

- **note** - A lightweight CLI tool for capturing quick notes and managing todos
EOF

# Initialize and push
git add .
git commit -m "feat: initial tap repository"
git push origin main
```

### 2. Create a GitHub Personal Access Token

GoReleaser needs a token to push to the Homebrew tap repository.

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Expiration: Choose based on your preference (recommend: No expiration)
5. Select scopes:
   - ✅ `repo` (Full control of private repositories)
   - ✅ `workflow` (Update GitHub Action workflows)
6. Click "Generate token"
7. **IMPORTANT**: Copy the token immediately (you won't be able to see it again)

### 3. Add the Token to Your Repository Secrets

1. Go to your `note` repository on GitHub
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
5. Value: Paste the token you just created
6. Click "Add secret"

### 4. Test GoReleaser Locally (Optional)

Before creating a release, you can test GoReleaser locally:

```bash
# Install GoReleaser
brew install goreleaser

# Test the configuration (dry run)
goreleaser release --snapshot --clean --skip=publish

# This will:
# - Build all binaries
# - Create archives
# - Generate checksums
# - But NOT publish anything
```

Check the `dist/` directory to see the generated artifacts.

### 5. Create Your First Release

Once everything is set up, create a release:

```bash
# Make changes and commit with conventional commits
git add .
git commit -m "feat: add new feature"
git push origin main

# Release-please will create a PR
# Review and merge the PR

# GoReleaser will automatically:
# ✅ Build binaries for all platforms
# ✅ Create GitHub release with binaries
# ✅ Update Homebrew tap with new formula
# ✅ Generate checksums
```

### 6. Verify Homebrew Installation

After the release is published, test the Homebrew installation:

```bash
# Add the tap
brew tap nathan-nicholson/tap

# Install note
brew install note

# Verify installation
note --help
```

## How It Works

### Release Flow

```
1. Commit with conventional commits
   └─> Push to main
       └─> release-please creates/updates PR
           └─> Merge PR
               └─> GitHub release created
                   └─> GoReleaser triggered
                       ├─> Build binaries (Linux, macOS, Windows)
                       ├─> Upload to GitHub release
                       └─> Update Homebrew tap
```

### Homebrew Formula Update

When GoReleaser runs, it:

1. Builds the project for macOS (amd64 and arm64)
2. Creates `.tar.gz` archives with checksums
3. Generates a Homebrew formula in `homebrew-tap/Formula/note.rb`
4. Commits and pushes to the tap repository
5. Users can then install with `brew install nathan-nicholson/tap/note`

### GoReleaser Configuration

The `.goreleaser.yml` file configures:

- **Builds**: Multi-platform compilation (Linux, macOS, Windows)
- **Archives**: TAR.GZ for Unix, ZIP for Windows
- **Checksums**: SHA256 for all artifacts
- **Homebrew**: Formula generation and tap publishing
- **Changelog**: Grouped by feature/fix/perf
- **Release**: GitHub release with custom header/footer

## Troubleshooting

### Issue: GoReleaser can't push to tap repository

**Solution**: Verify the `HOMEBREW_TAP_GITHUB_TOKEN` secret:
- Has `repo` and `workflow` scopes
- Is not expired
- Is correctly named in repository secrets

### Issue: CGO compilation fails

**Solution**: The workflow installs cross-compilation tools. For local testing:
```bash
# macOS
brew install FiloSottile/musl-cross/musl-cross

# Linux
sudo apt-get install gcc-aarch64-linux-gnu
```

### Issue: Formula not found after release

**Solution**:
1. Check that the release completed successfully
2. Verify the tap repository has the formula in `Formula/note.rb`
3. Update the tap: `brew update`
4. Try untapping and retapping: `brew untap nathan-nicholson/tap && brew tap nathan-nicholson/tap`

## Advanced Configuration

### Supporting Older macOS Versions

To support older macOS versions, update `.goreleaser.yml`:

```yaml
builds:
  - env:
      - CGO_ENABLED=1
      - MACOSX_DEPLOYMENT_TARGET=10.15  # Support macOS 10.15+
```

### Adding Bottles (Precompiled Binaries)

Homebrew bottles provide faster installation. To add bottle support:

1. After the formula is created, build bottles for different macOS versions
2. Use `brew test-bot` to build and upload bottles
3. This is advanced and typically done by the Homebrew core team

For now, the formula builds from source, which is fine for most users.

### Custom Formula Options

Edit `.goreleaser.yml` to customize the Homebrew formula:

```yaml
brews:
  - name: note
    caveats: |
      Data is stored in ~/.note/notes.db

      To get started:
        note "My first note"
        note todo "My first task"

    conflicts:
      - other-note-tool

    dependencies:
      - name: sqlite
        type: run
```

## Resources

- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Release-Please Documentation](https://github.com/googleapis/release-please)
- [Conventional Commits](https://www.conventionalcommits.org/)

## Support

If you encounter issues:

1. Check the GitHub Actions logs for errors
2. Verify all secrets are correctly configured
3. Test GoReleaser locally with `--snapshot --skip=publish`
4. Review the Homebrew tap repository for the generated formula
