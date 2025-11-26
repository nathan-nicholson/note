# CLAUDE.md - AI Assistant Guide for `note`

> **Purpose**: This file serves as Claude Code's memory for the `note` project. It provides comprehensive context about the codebase structure, architecture, development practices, and conventions. Read this file carefully before making any code changes.

## How to Use This Documentation

**For Claude Code**:
- **Before coding**: Review relevant sections to understand existing patterns
- **When adding features**: Follow the "Common Tasks" section for step-by-step guidance
- **When unsure**: Check "Code Conventions" and "Best Practices Summary"
- **Before committing**: Review "Questions to Ask Before Making Changes"
- **When stuck**: Consult the "Troubleshooting Guide"

**When to ask for clarification vs. proceed**:
- ✅ **Proceed autonomously** when task matches documented patterns (adding commands, fields, tests)
- ✅ **Proceed autonomously** for bug fixes with clear root causes
- ❓ **Ask for clarification** when user intent is ambiguous (multiple valid approaches)
- ❓ **Ask for clarification** when changes would significantly alter architecture
- ❓ **Ask for clarification** when breaking changes are required

## Project Overview

`note` is a lightweight CLI tool written in Go for capturing quick notes and managing todos with project-based organization. It provides a keyboard-driven interface for thought capture, task tracking, and project management.

**Key Features:**
- Quick note capture from terminal
- Todo management with due dates
- Project-based organization with automatic tagging
- Activity logging for lifecycle events
- SQLite-based local storage (`~/.note/notes.db`)
- Cross-platform support (Linux, macOS, Windows)

**Tech Stack:**
- Language: Go 1.24+
- Database: SQLite3 (with CGO)
- CLI Framework: Cobra
- Display: fatih/color for terminal output
- Build: GoReleaser for multi-platform releases
- CI/CD: GitHub Actions with release-please

### What Makes This Project Unique

**Design Philosophy**:
- **Simplicity First**: Minimal cognitive overhead for users
- **Local-First**: All data stored locally in SQLite, no cloud dependencies
- **Activity Awareness**: Automatic logging of significant events for context
- **Project Context**: Automatic tagging based on active project eliminates manual organization

**Key Architectural Decisions**:
- **Global DB Connection**: Simplifies code but requires proper initialization order
- **CGO Dependency**: SQLite requires CGO, complicating cross-compilation
- **Migrations on Startup**: Schema changes applied automatically, no separate migration tool
- **Natural Language Dates**: User-friendly input reduces friction

**Notable Constraints**:
- Single active project at a time (singleton pattern)
- Project names must be kebab-case (enforced validation)
- Migrations are append-only (never modify existing)
- Activity notes created automatically (cannot be disabled)

## Repository Structure

```
note/
├── cmd/                    # CLI command definitions (Cobra)
│   ├── root.go            # Root command and CLI setup
│   ├── add.go             # Note creation
│   ├── list.go            # Note listing
│   ├── edit.go            # Note editing
│   ├── delete.go          # Note deletion
│   ├── show.go            # Note display
│   ├── tags.go            # Tag management
│   ├── todo*.go           # Todo commands (add, list, edit, complete, etc.)
│   ├── project*.go        # Project commands (create, close, status, etc.)
│   ├── version.go         # Version command
│   └── update.go          # Update checker
├── internal/              # Private application code
│   ├── activity/          # Activity logging (todo/project events)
│   ├── database/          # DB initialization and migrations
│   ├── dateparse/         # Natural language date parsing
│   ├── display/           # Terminal output formatting
│   ├── models/            # Data structures (Note, Todo, Project, Tag)
│   ├── repository/        # Data access layer (CRUD operations)
│   ├── update/            # Version checking and updates
│   └── version/           # Version information
├── main.go                # Application entry point
├── Makefile               # Build and test commands
├── .goreleaser.yml        # Multi-platform release configuration
├── .github/workflows/     # CI/CD pipelines
│   ├── ci.yml            # Test, build, lint
│   ├── release-please.yml # Automated release PR creation
│   └── release.yml        # Binary builds on release
└── README.md              # User documentation
```

## Architecture Patterns

### Layered Architecture

The codebase follows a clean layered architecture:

1. **Command Layer** (`cmd/`): CLI interface, argument parsing, user interaction
2. **Repository Layer** (`internal/repository/`): Database operations, CRUD
3. **Model Layer** (`internal/models/`): Data structures and validation
4. **Database Layer** (`internal/database/`): Connection management, migrations
5. **Support Services**: Display formatting, date parsing, activity logging

### Data Flow

```
User Input (CLI)
    ↓
Cobra Command (cmd/)
    ↓
Repository Function (internal/repository/)
    ↓
SQL Query Execution
    ↓
Database (SQLite)
    ↓
Model Struct (internal/models/)
    ↓
Display Formatter (internal/display/)
    ↓
Terminal Output
```

### Key Design Principles

- **Global DB Connection**: Single `database.DB` global variable initialized at startup
- **Activity Logging**: Automatic note creation for significant todo/project events
- **Project Auto-Tagging**: Active project automatically tags new notes and todos
- **Natural Language Dates**: Parser supports "tomorrow", "end-of-week", ISO format
- **Validation**: Model-level validation (e.g., project name kebab-case)

## Core Components

### 1. Database Layer (`internal/database/`)

**Files**: `db.go`, `migrations.go`

**Responsibilities**:
- Initialize SQLite connection at `~/.note/notes.db`
- Run schema migrations on startup
- Ensure "home" default project exists
- Provide global `DB` variable for repository access

**Key Functions**:
- `InitDB()`: Create DB directory, open connection, run migrations
- `CloseDB()`: Cleanup on application exit
- `runMigrations()`: Execute schema updates
- `ensureHomeProject()`: Create default project if missing

### 2. Models (`internal/models/`)

**Files**: `note.go`, `todo.go`, `project.go`, `tag.go`

**Data Structures**:

```go
type Note struct {
    ID          int
    Content     string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    IsImportant bool
    Tags        []string
}

type Todo struct {
    ID          int
    Content     string
    IsComplete  bool
    DueDate     sql.NullTime
    CreatedAt   time.Time
    UpdatedAt   time.Time
    CompletedAt sql.NullTime
    Tags        []string
}

type Project struct {
    ID               int
    Name             string
    CreatedAt        time.Time
    FirstActivatedAt sql.NullTime
    LastActivityAt   sql.NullTime
    ClosedAt         sql.NullTime
    IsClosed         bool
    Tags             []string
}
```

**Validation**:
- Project names must be kebab-case (`^[a-z0-9]+(-[a-z0-9]+)*$`)
- Reserved project names: create, close, reopen, list, status, show, edit, delete

### 3. Repository Layer (`internal/repository/`)

**Files**: `note_repo.go`, `todo_repo.go`, `project_repo.go`, `tag_repo.go`

**Pattern**: All repository functions accept `*sql.DB` as first parameter

**Key Functions by Domain**:

**Notes**:
- `CreateNote(db, content, tags, isImportant) (*Note, error)`
- `GetNotes(db, startDate, endDate, tags, important, showActivityNotes) ([]*Note, error)`
- `UpdateNote(db, id, content, tags, isImportant) error`
- `DeleteNote(db, id) error`

**Todos**:
- `CreateTodo(db, content, tags, dueDate) (*Todo, error)`
- `GetTodos(db, incomplete, tags) ([]*Todo, error)`
- `CompleteTodo(db, id) error`
- `UncompleteTodo(db, id) error`

**Projects**:
- `CreateProject(db, name, tags) (*Project, error)`
- `GetActiveProject(db) (*Project, error)`
- `ActivateProject(db, name) error`
- `CloseProject(db, name) error` (requires all todos complete)

### 4. Activity Logger (`internal/activity/`)

**Purpose**: Automatically create notes for important lifecycle events

**Logged Events**:
- Todo: create, update, complete, delete
- Project: create, activate, deactivate, update, close, reopen, delete

**Pattern**: All log functions create notes with specific tags
- Example: `LogTodoCompleted()` creates note with tags `["todo", "complete", ...todo.Tags]`

### 5. Date Parser (`internal/dateparse/`)

**Natural Language Support**:
- `today`: Current date
- `tomorrow`: Next day
- `end-of-week`: Next Friday
- `end-of-month`: Last day of current month
- `next-week`: 7 days from today
- `next-month`: 30 days from today
- ISO: `2025-11-25`

**Function**: `ParseDate(input string) (time.Time, error)`

### 6. Display Formatters (`internal/display/`)

**Files**: `formatter.go`, `note_list.go`, `todo_list.go`, `project_list.go`

**Responsibilities**:
- Format notes with timestamps, tags, importance markers
- Group todos by status (Overdue, Today, This Week, Future, No Due Date, Completed)
- Display project status with todo summaries
- Color-coded terminal output using `fatih/color`

## Development Workflows

### Setting Up Development Environment

```bash
# Clone repository
git clone https://github.com/nathan-nicholson/note.git
cd note

# Install dependencies
go mod download

# Build binary
make build

# Run tests
make test

# Install locally
make install
```

### Building and Testing

**Makefile targets**:
- `make build`: Compile binary to `./note`
- `make install`: Install to `$GOPATH/bin`
- `make test`: Run all tests
- `make test-verbose`: Tests with verbose output
- `make test-coverage`: Generate HTML coverage report
- `make clean`: Remove build artifacts

**Manual commands**:
```bash
go build -o note main.go
go test ./...
go test -v -race -coverprofile=coverage.out ./...
go test ./internal/dateparse/...
```

### Running the Application

```bash
# Create a note
./note "Meeting notes"

# Create todo
./note todo "Review PR" --due tomorrow

# Switch project
./note project work

# List notes
./note list --start 2025-11-20
```

## Testing Practices

### Test Structure

- **Unit tests**: `*_test.go` files alongside source
- **Integration tests**: Repository tests with in-memory SQLite
- **Test utilities**: `testutil_test.go` provides test database setup

### Current Coverage

- Date parser: ~95%
- Repository layer: ~48% (focus area for improvement)
- Activity logger: Needs more coverage
- Display formatters: Minimal coverage

### Writing Tests

**Pattern for repository tests**:

```go
func TestSomething(t *testing.T) {
    db := setupTestDB(t) // Creates in-memory DB with migrations
    defer db.Close()

    // Test logic
    result, err := repository.SomeFunction(db, args...)

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Assertions
}
```

**Test database setup** (`internal/repository/testutil_test.go`):
- Creates in-memory SQLite instance
- Runs all migrations
- Ensures clean state per test

### Testing Strategy and Priorities

**What to test thoroughly** (high value):
1. **Repository layer** - All CRUD operations, edge cases, SQL correctness
2. **Model validation** - Project name rules, data constraints
3. **Date parsing** - All natural language formats, edge cases
4. **Activity logging** - Verify notes created for events
5. **Project auto-tagging** - Tags applied correctly

**What to test lightly** (medium value):
1. **Display formatters** - Sample output verification
2. **Command parsing** - Flag handling, argument validation
3. **Error messages** - User-facing error text

**What can be tested manually** (lower priority for automation):
1. **Terminal color output** - Visual inspection
2. **Update checker** - External API integration
3. **Cross-platform builds** - Handled by CI

**Test Coverage Goals**:
- Repository: Target 80%+ (currently ~48%)
- Date parser: Maintain 95%+
- Models: Target 90%+
- Activity logger: Target 80%+

**When adding new features**:
- Write tests BEFORE or DURING implementation, not after
- Test happy path AND error cases
- Test edge cases (empty strings, null values, boundaries)
- Add integration tests for multi-step workflows

## CI/CD Pipeline

### Continuous Integration (`.github/workflows/ci.yml`)

**Triggers**: Push and PR to `main`

**Jobs**:
1. **Test**: Run tests on Go 1.22, 1.23, 1.24 with race detector
2. **Build**: Compile binary and verify functionality
3. **Lint**: Run golangci-lint

**Coverage**: Reports uploaded to Codecov

### Release Process (release-please + GoReleaser)

**Automated Release Flow**:

1. **Development**: Commit with conventional commit format
   ```bash
   feat: add recurring todos
   fix: correct date parsing for leap years
   docs: update README
   ```

2. **Release PR Creation** (`.github/workflows/release-please.yml`):
   - release-please monitors commits
   - Creates/updates release PR automatically
   - Generates CHANGELOG from conventional commits
   - Bumps version (feat=minor, fix=patch, BREAKING CHANGE=major)

3. **Release Trigger** (`.github/workflows/release.yml`):
   - Merge release PR
   - release-please creates GitHub release with tag
   - GoReleaser workflow triggered

4. **Binary Building** (`.goreleaser.yml`):
   - Cross-compile for Linux (amd64, arm64), macOS (Intel, Apple Silicon), Windows (amd64)
   - CGO enabled for Unix (SQLite requirement)
   - CGO disabled for Windows (simpler cross-compilation)
   - Generate checksums
   - Upload artifacts to release
   - Update Homebrew tap automatically

### Conventional Commit Types

- `feat`: New feature (minor version bump)
- `fix`: Bug fix (patch version bump)
- `docs`: Documentation only
- `test`: Adding/updating tests
- `perf`: Performance improvements
- `refactor`: Code refactoring
- `build`: Build system changes
- `ci`: CI/CD changes
- `chore`: Maintenance tasks

**Breaking changes**: Add `!` or `BREAKING CHANGE:` footer for major version bump

### GoReleaser Configuration Highlights

**CGO Cross-Compilation**:
- macOS amd64: Uses `o64-clang`
- macOS arm64: Uses `oa64-clang`
- Linux amd64: Uses `gcc`
- Linux arm64: Uses `aarch64-linux-gnu-gcc`
- Windows: CGO disabled (no SQLite dependency issues)

**Homebrew Publishing**:
- Automatically updates `nathan-nicholson/homebrew-tap`
- Requires `HOMEBREW_TAP_GITHUB_TOKEN` secret
- Includes SQLite dependency

**Version Injection**:
```go
-ldflags: -X github.com/nathan-nicholson/note/internal/version.Version={{.Version}}
```

## Code Conventions for AI Assistants

### 1. File Modification Guidelines

**ALWAYS**:
- Read existing files before modifying
- Preserve existing code style and patterns
- Use the repository layer for all database operations
- Include error handling for all operations
- Update tests when changing functionality
- Run `make test` before committing
- Use `gofmt` formatting (automatic with most editors)
- Follow Go naming conventions (PascalCase for exports, camelCase for private)

**NEVER**:
- Modify `main.go` unless absolutely necessary
- Change database migration logic without careful review
- Skip validation in model layer
- Add dependencies without justification
- Create new global variables
- Commit code that doesn't pass tests
- Use `panic()` except in `init()` or truly unrecoverable situations
- Ignore golangci-lint warnings without good reason

### Code Quality Standards

**Go Best Practices**:
- Use meaningful variable names (avoid single letters except loop counters)
- Keep functions focused and small (prefer <50 lines)
- Return errors, don't panic
- Use `defer` for cleanup (database connections, file handles)
- Avoid naked returns in functions >10 lines
- Use `context.Context` for cancellation in future additions

**This Codebase Specifically**:
- Database operations always use prepared statements (SQL injection prevention)
- All user-facing strings go through display formatters
- Error messages should be actionable (tell user what to do)
- Commands should work with or without flags (sensible defaults)
- Prefer explicit over implicit (don't hide important operations)

### 2. Adding New Commands

**Pattern to follow** (see `cmd/` examples):

```go
package cmd

import (
    "github.com/nathan-nicholson/note/internal/database"
    "github.com/nathan-nicholson/note/internal/repository"
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new [args]",
    Short: "Brief description",
    Long:  `Detailed description`,
    Args:  cobra.ExactArgs(1), // Or other validator
    RunE: func(cmd *cobra.Command, args []string) error {
        // 1. Parse arguments
        // 2. Call repository functions
        // 3. Use database.DB global variable
        // 4. Handle errors
        // 5. Return nil on success
        return nil
    },
}

func init() {
    // Add flags
    newCmd.Flags().StringVar(&variable, "flag", "", "description")

    // Register with root in root.go:
    // rootCmd.AddCommand(newCmd)
}
```

### 3. Adding Repository Functions

**Pattern** (`internal/repository/`):

```go
func GetSomething(db *sql.DB, filter string) (*models.Thing, error) {
    query := `SELECT id, field1, field2 FROM table WHERE field = ?`

    row := db.QueryRow(query, filter)

    var thing models.Thing
    err := row.Scan(&thing.ID, &thing.Field1, &thing.Field2)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("thing not found")
    }
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &thing, nil
}
```

**Multi-row queries** use `db.Query()` and iterate:

```go
rows, err := db.Query(query, args...)
if err != nil {
    return nil, err
}
defer rows.Close()

var results []*models.Thing
for rows.Next() {
    var thing models.Thing
    if err := rows.Scan(&thing.ID, &thing.Field); err != nil {
        return nil, err
    }
    results = append(results, &thing)
}
return results, rows.Err()
```

### 4. Database Migrations

**Location**: `internal/database/migrations.go`

**Adding new migration**:

```go
func runMigrations(db *sql.DB) error {
    migrations := []string{
        // Existing migrations...
        `CREATE TABLE IF NOT EXISTS new_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            field TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
    }

    for _, migration := range migrations {
        if _, err := db.Exec(migration); err != nil {
            return err
        }
    }
    return nil
}
```

**CRITICAL**: Migrations are append-only. Never modify existing migrations.

### 5. Activity Logging

When adding new operations that should be logged:

```go
// In internal/activity/logger.go
func LogNewActivity(db *sql.DB, details string) error {
    content := fmt.Sprintf("Activity: %s", details)
    tags := []string{"activity-type", "action"}
    _, err := repository.CreateNote(db, content, tags, false)
    return err
}

// In command handler
if err := activity.LogNewActivity(database.DB, details); err != nil {
    // Log error but don't fail operation
}
```

### 6. Error Handling Patterns

**Repository layer**:
```go
if err == sql.ErrNoRows {
    return nil, fmt.Errorf("resource not found")
}
if err != nil {
    return nil, fmt.Errorf("database error: %w", err)
}
```

**Command layer**:
```go
result, err := repository.DoSomething(database.DB, args)
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

**Main**: Errors propagate to `main.go` which prints to stderr and exits

### 7. Tag Handling

**Tags are always string slices**:
- Automatically lowercased during processing
- Project auto-tagging: append active project name to tags
- Activity notes: prepend event type tags

```go
// Get active project
activeProject, _ := repository.GetActiveProject(database.DB)

// Combine tags
tags := append(userTags, activeProject.Name)
```

### 8. Display Formatting

**Use existing display functions**:
```go
import "github.com/nathan-nicholson/note/internal/display"

notes, _ := repository.GetNotes(...)
display.FormatNoteList(notes, showIDs)
```

**Color usage**:
```go
import "github.com/fatih/color"

color.Green("Success message")
color.Yellow("Warning message")
color.Red("Error message")
color.Cyan("Info message")
```

### 9. Testing New Code

**Always add tests for**:
- New repository functions
- Date parsing additions
- Model validation logic
- Complex business logic

**Use test utilities**:
```go
func TestNewFeature(t *testing.T) {
    db := setupTestDB(t) // Provided by testutil_test.go
    defer db.Close()

    // Test implementation
}
```

### 10. Version Updates

**Version controlled by release-please**:
- `internal/version/version.go` contains version constant
- Updated automatically by release-please
- Injected at build time via ldflags in GoReleaser

**Manual version check**:
```bash
note version
```

## Common Scenarios and How to Handle Them

### Scenario: User asks "Add a search command"

**Your approach**:
1. Review existing commands in `cmd/` to understand patterns
2. Check if repository layer has search capabilities (it doesn't - need to add)
3. Follow "Task: Add a New Command" section below
4. Add repository function for search
5. Update display formatter for search results
6. Add tests for search functionality
7. Update README with search examples
8. Commit with `feat: add search command for notes`

### Scenario: User reports "Date parsing doesn't work for 'next friday'"

**Your approach**:
1. Check `internal/dateparse/parser.go` to see if "next friday" is supported
2. Review existing test cases in `internal/dateparse/parser_test.go`
3. Add test case that reproduces the issue
4. Implement the new date format
5. Verify test passes
6. Check if README needs updating
7. Commit with `fix: add support for 'next friday' date parsing`

### Scenario: User asks "Why is the database slow?"

**Your approach**:
1. Ask clarifying questions about the specific operation that's slow
2. Review relevant repository function in `internal/repository/`
3. Check if indexes exist in `internal/database/migrations.go`
4. Suggest adding indexes if missing
5. Follow migration pattern to add index
6. Test performance improvement
7. Commit with `perf: add index to notes table for faster queries`

### Scenario: User requests "Export notes to JSON"

**Your approach**:
1. Check if this fits existing command pattern or needs new command
2. Decide: Should this be `note export` or `note list --format json`?
3. Ask user for preference if ambiguous
4. Implement based on decision
5. Add tests for JSON export
6. Update README
7. Commit with `feat: add JSON export for notes`

## Common Tasks for AI Assistants

### Task: Add a New Note Field

1. **Update model** (`internal/models/note.go`):
   ```go
   type Note struct {
       // ... existing fields
       NewField string
   }
   ```

2. **Add migration** (`internal/database/migrations.go`):
   ```go
   `ALTER TABLE notes ADD COLUMN new_field TEXT`
   ```

3. **Update repository functions** (`internal/repository/note_repo.go`):
   - Modify `CreateNote()` to accept new field
   - Update SQL queries to include new column
   - Update `Scan()` calls

4. **Update commands** (`cmd/add.go`, `cmd/edit.go`):
   - Add flag for new field
   - Pass to repository functions

5. **Update display** (`internal/display/note_list.go`):
   - Include new field in output

6. **Add tests** (`internal/repository/note_repo_test.go`):
   - Test creation with new field
   - Test querying with new field

### Task: Add a New Command

1. **Create command file** (`cmd/newcmd.go`):
   ```go
   package cmd

   var newCmd = &cobra.Command{
       Use:   "newcmd",
       Short: "Description",
       RunE: func(cmd *cobra.Command, args []string) error {
           // Implementation
           return nil
       },
   }
   ```

2. **Register in root** (`cmd/root.go`):
   ```go
   func init() {
       rootCmd.AddCommand(newCmd)
   }
   ```

3. **Add repository function** if needed (`internal/repository/`)

4. **Update README** with usage examples

5. **Add tests** for new functionality

### Task: Improve Test Coverage

1. **Identify untested code**:
   ```bash
   make test-coverage
   open coverage.html
   ```

2. **Focus areas** (current low coverage):
   - Repository layer (~48%)
   - Activity logger
   - Display formatters
   - Update checker

3. **Write tests** following existing patterns in `*_test.go` files

4. **Verify coverage improvement**:
   ```bash
   make test-coverage
   ```

### Task: Fix a Bug

1. **Reproduce the issue**: Create test case that fails
2. **Identify root cause**: Check relevant repository/command code
3. **Implement fix**: Modify minimal code necessary
4. **Verify fix**: Ensure test passes
5. **Check side effects**: Run full test suite
6. **Commit with conventional format**:
   ```bash
   git commit -m "fix: describe what was fixed"
   ```

### Task: Add Date Format

1. **Update parser** (`internal/dateparse/parser.go`):
   ```go
   func ParseDate(input string) (time.Time, error) {
       // Add new case
       case "new-format":
           return calculateDate(), nil
   ```

2. **Add tests** (`internal/dateparse/parser_test.go`):
   ```go
   func TestParseDate_NewFormat(t *testing.T) {
       result, err := ParseDate("new-format")
       // Assertions
   }
   ```

3. **Update documentation** (README.md, help text)

## Database Schema Overview

### Tables

**notes**:
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `content`: TEXT NOT NULL
- `created_at`: DATETIME DEFAULT CURRENT_TIMESTAMP
- `updated_at`: DATETIME DEFAULT CURRENT_TIMESTAMP
- `is_important`: BOOLEAN DEFAULT 0

**todos**:
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `content`: TEXT NOT NULL
- `is_complete`: BOOLEAN DEFAULT 0
- `due_date`: DATE NULL
- `created_at`: DATETIME DEFAULT CURRENT_TIMESTAMP
- `updated_at`: DATETIME DEFAULT CURRENT_TIMESTAMP
- `completed_at`: DATETIME NULL

**projects**:
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `name`: TEXT UNIQUE NOT NULL
- `created_at`: DATETIME DEFAULT CURRENT_TIMESTAMP
- `first_activated_at`: DATETIME NULL
- `last_activity_at`: DATETIME NULL
- `closed_at`: DATETIME NULL
- `is_closed`: BOOLEAN DEFAULT 0

**tags**: (association tables)
- `note_tags`: note_id, tag
- `todo_tags`: todo_id, tag
- `project_tags`: project_id, tag

**active_project**:
- `id`: Always 1 (singleton)
- `project_id`: Foreign key to projects

## Troubleshooting Guide

### CGO/SQLite Build Issues

**Problem**: `undefined: sqlite3` or CGO errors

**Solution**:
- Ensure CGO_ENABLED=1 for Unix builds
- Install sqlite3 development libraries
- For cross-compilation, use appropriate C compiler (see `.goreleaser.yml`)

### Migration Failures

**Problem**: Database migration errors on startup

**Solution**:
- Check `internal/database/migrations.go` for syntax errors
- Never modify existing migrations
- For development, delete `~/.note/notes.db` and restart
- For production, add compensating migration

### Test Database Issues

**Problem**: Tests failing with database errors

**Solution**:
- Ensure `setupTestDB()` is called in each test
- Check that database is closed with `defer db.Close()`
- Verify migrations run in test setup

### Release Pipeline Issues

**Problem**: GoReleaser fails on release

**Solution**:
- Check `HOMEBREW_TAP_GITHUB_TOKEN` secret is set
- Verify `.goreleaser.yml` syntax
- Ensure cross-compilation tools are available in CI
- Check release-please created valid tag

## Key Files Reference

| File | Purpose |
|------|---------|
| `main.go` | Entry point, DB initialization |
| `cmd/root.go` | Root command, subcommand registration |
| `internal/database/db.go` | Database connection management |
| `internal/database/migrations.go` | Schema definitions |
| `internal/repository/*_repo.go` | Data access layer |
| `internal/models/*.go` | Data structures |
| `internal/activity/logger.go` | Event logging |
| `internal/dateparse/parser.go` | Natural language date parsing |
| `internal/display/*.go` | Terminal output formatting |
| `Makefile` | Build and test commands |
| `.goreleaser.yml` | Multi-platform release config |
| `.github/workflows/ci.yml` | CI pipeline |
| `.github/workflows/release-please.yml` | Automated releases |

## Best Practices Summary

1. **Always read before modifying** - Understand existing code patterns
2. **Use repository layer** - Never write SQL in command layer
3. **Follow error handling patterns** - Wrap errors with context
4. **Write tests** - Especially for repository functions
5. **Use conventional commits** - Enable automated releases
6. **Preserve existing style** - Match surrounding code
7. **Update documentation** - Keep README and help text current
8. **Test locally** - Run `make test` before committing
9. **Check coverage** - Use `make test-coverage` to identify gaps
10. **Activity logging** - Log significant user-facing events

## Questions to Ask Before Making Changes

1. Does this change require a database migration?
2. Should this operation create an activity log entry?
3. Are there existing patterns I should follow?
4. Do I need to update the display formatter?
5. Does this affect project auto-tagging?
6. Should this be a new command or a flag on existing command?
7. What error cases should I handle?
8. What tests should I add?
9. Does this change the user-facing API?
10. Should this be documented in README?

## File Organization and Imports

### CLAUDE.md File Hierarchy

This file is located at the repository root and provides project-wide context. Claude Code loads CLAUDE.md files recursively from the current directory up to the repository root.

**For subdirectory-specific context**:
- Create `subdirectory/CLAUDE.md` for component-specific guidance
- More specific files build upon or refine base knowledge
- Example: `internal/repository/CLAUDE.md` could document repository-specific patterns

**For local-only context** (not committed to git):
- Use `CLAUDE.local.md` for machine-specific settings
- Add to `.gitignore` to prevent committing

### Importing Additional Documentation

This file can import other documentation using `@path/to/file` syntax:

```markdown
## Additional Context

@./CONTRIBUTING.md
@./docs/architecture.md
```

**Note**: Currently, all relevant context is self-contained in this file. Use imports if the file grows beyond ~1000 lines or if specific components need dedicated documentation.

---

**Last Updated**: 2025-11-26
**Codebase Version**: 1.0.0
**Total Lines of Go Code**: ~1,268

**Verification**: Use `/memory` command in Claude Code to verify this file is loaded.
