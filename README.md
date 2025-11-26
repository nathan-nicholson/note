# note

A lightweight CLI tool for capturing quick notes and managing todos with project-based organization.

## Installation

### Homebrew (Recommended)

```bash
brew tap nathan-nicholson/tap
brew install note
```

### Go Install

```bash
go install github.com/nathan-nicholson/note@latest
```

### Download Binary

Download the latest release for your platform from the [Releases](https://github.com/nathan-nicholson/note/releases) page.

Available for:
- Linux (amd64, arm64)
- macOS (Intel, Apple Silicon)
- Windows (amd64)

### Build from Source

```bash
git clone https://github.com/nathan-nicholson/note.git
cd note
make build
./note --help
```

## Quick Start

```bash
note "Quick thought"
note todo "Task to complete" --due tomorrow
note project work
```

## Features

- **Quick note capture** - Instantly save thoughts from the terminal
- **Todo management** - Track tasks with due dates and completion status
- **Project-based organization** - Group work by projects with automatic tagging
- **Activity logging** - Automatic notes for todo and project lifecycle events
- **Flexible tagging** - Tag notes and todos for easy filtering
- **Date-based filtering** - Find notes by date range

## Usage

### Notes

Create a note:
```bash
note "Meeting with team at 2pm"
note "Important decision" --important --tag architecture
```

List notes:
```bash
note list                                    # Today's notes
note list --start 2025-11-20                 # From Nov 20 to today
note list --start 2025-11-20 --end 2025-11-22
note list --tag work --important             # Filter by tags and importance
note list --show-ids                         # Show IDs for editing
```

Edit and manage:
```bash
note edit 42 --content "Updated content"
note edit 42 --tag newTag
note show 42
note delete 42
```

### Todos

Create todos:
```bash
note todo "Review PR"
note todo "File taxes" --due 2025-12-31 --tag finance
note todo "Weekly report" --due tomorrow
```

List todos:
```bash
note todo list                               # All todos grouped by status
note todo list --incomplete                  # Only incomplete
note todo list --tag work                    # Filter by tag
```

Manage todos:
```bash
note todo complete 42
note todo uncomplete 42
note todo edit 42 --content "Updated task" --due next-week
note todo show 42
note todo delete 42
```

### Projects

Create and switch projects:
```bash
note project create work --tag professional
note project work                            # Switch to work project
```

List projects:
```bash
note project list                            # Open projects only
note project list --all                      # Include closed projects
```

Project status:
```bash
note project status                          # Current project
note project status work                     # Specific project
note project status --all                    # Show all tasks including completed
```

Close and manage:
```bash
note project close work                      # Close project (all todos must be complete)
note project reopen work                     # Reopen a closed project
note project edit work --tag professional --tag fulltime
note project show work
note project delete old-project
```

### Tags

List all tags with usage counts:
```bash
note tags
```

### Version & Updates

Check current version:
```bash
note version
```

Check for updates:
```bash
note update --check              # Check for updates only
note update                      # Check and prompt to install
note update -y                   # Check and install without prompt
```

The update command will:
- Check for the latest version on GitHub
- Show release notes for the new version
- Detect your installation method (Homebrew, go install, or binary)
- Automatically update if possible (Homebrew/go install)
- Provide download instructions for manual installs

## Date Formats

Natural language:
- `today`
- `tomorrow`
- `end-of-week` (next Friday)
- `end-of-month`
- `next-week` (7 days from today)
- `next-month`

ISO format:
- `2025-11-25`
- `2025-12-31`

## Data Storage

All data is stored in `~/.note/notes.db` using SQLite.

## Project Auto-Tagging

When a project is active, all new notes and todos are automatically tagged with the project name:

```bash
note project work
note "Sprint planning complete"
# Automatically tagged: #work

note todo "Deploy to production" --due friday
# Automatically tagged: #work
```

## Activity Notes

The tool automatically creates notes for important events:

**Todo activities:**
- Created todo
- Updated todo
- Completed todo
- Deleted todo

**Project activities:**
- Created project
- Activated project
- Deactivated project
- Updated project
- Closed project
- Reopened project
- Deleted project

These notes appear in your regular note list and can be filtered using tags like `#todo`, `#project`, `#complete`, etc.

## Development

### Building and Testing

```bash
make build           # Build the binary
make install         # Install to $GOPATH/bin
make test            # Run tests
make test-verbose    # Run tests with verbose output
make test-coverage   # Generate HTML coverage report
make clean           # Remove build artifacts
```

### Running Tests

The project has comprehensive unit and integration tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/dateparse/...
go test ./internal/repository/...
```

Current test coverage:
- Date parser: 95.2%
- Repository layer: 47.5%

## CI/CD

This project uses GitHub Actions for continuous integration and automated releases.

### Continuous Integration

On every push and pull request to `main`:
- Tests run on Go 1.22, 1.23, and 1.24
- Code is built and verified
- Linting checks are performed
- Coverage reports are generated

### Automated Releases

This project uses [release-please](https://github.com/googleapis/release-please) for automated releases:

1. **Commit with Conventional Commits** - Use conventional commit messages (see below)
2. **Release PR Created** - Release-please creates/updates a release PR
3. **Merge Release PR** - When merged, a new release is created automatically
4. **Binaries Built** - Cross-platform binaries are built and attached to the release

### Conventional Commits

Use conventional commit format for your commits:

```bash
# Features
git commit -m "feat: add support for recurring todos"
git commit -m "feat(projects): add project archiving"

# Bug fixes
git commit -m "fix: correct date parsing for leap years"
git commit -m "fix(database): handle connection timeout"

# Documentation
git commit -m "docs: update installation instructions"

# Tests
git commit -m "test: add integration tests for projects"

# Performance improvements
git commit -m "perf: optimize note listing query"

# Refactoring
git commit -m "refactor: simplify tag management logic"

# CI/Build changes
git commit -m "ci: add code coverage reporting"
git commit -m "build: update dependencies"
```

**Types:**
- `feat`: New feature (triggers minor version bump)
- `fix`: Bug fix (triggers patch version bump)
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `perf`: Performance improvements
- `refactor`: Code refactoring
- `build`: Build system or dependency changes
- `ci`: CI/CD configuration changes
- `chore`: Other changes that don't modify src or test files

**Breaking Changes:**
For breaking changes, add `!` after the type or add `BREAKING CHANGE:` in the footer:

```bash
git commit -m "feat!: remove deprecated todo filters"

# Or with footer
git commit -m "feat: redesign project API

BREAKING CHANGE: Project API now uses new schema"
```

### Release Binaries

When a release is created using release-please and GoReleaser:
- **Binaries** are built for Linux (amd64, arm64), macOS (Intel, Apple Silicon), and Windows (amd64)
- **Homebrew formula** is automatically updated in [nathan-nicholson/homebrew-tap](https://github.com/nathan-nicholson/homebrew-tap)
- **Checksums** are generated for all artifacts
- **Release notes** are automatically generated from commits

Download the latest release from the [Releases](https://github.com/nathan-nicholson/note/releases) page.

**Setting up Homebrew publishing**: See [HOMEBREW_SETUP.md](HOMEBREW_SETUP.md) for detailed instructions.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes using conventional commits
4. Push to your branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

All PRs must:
- Pass CI checks (tests, linting)
- Include tests for new functionality
- Follow conventional commit format
- Update documentation as needed

## License

MIT
