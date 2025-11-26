package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/nathan-nicholson/note/internal/version"
)

const (
	githubAPIURL = "https://api.github.com/repos/nathan-nicholson/note/releases/latest"
	timeout      = 10 * time.Second
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Name    string `json:"name"`
	Body    string `json:"body"`
}

type UpdateInfo struct {
	CurrentVersion string
	LatestVersion  string
	UpdateURL      string
	IsNewer        bool
	ReleaseNotes   string
}

type InstallMethod int

const (
	InstallUnknown InstallMethod = iota
	InstallHomebrew
	InstallGoInstall
	InstallBinary
)

// CheckForUpdates checks if a newer version is available
func CheckForUpdates() (*UpdateInfo, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent to avoid GitHub API rate limiting
	req.Header.Set("User-Agent", "note-cli")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	currentVersion := normalizeVersion(version.Version)
	latestVersion := normalizeVersion(release.TagName)

	updateInfo := &UpdateInfo{
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		UpdateURL:      release.HTMLURL,
		IsNewer:        compareVersions(latestVersion, currentVersion) > 0,
		ReleaseNotes:   release.Body,
	}

	return updateInfo, nil
}

// DetectInstallMethod tries to determine how note was installed
func DetectInstallMethod() InstallMethod {
	// Check if installed via Homebrew
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if isHomebrewInstall() {
			return InstallHomebrew
		}
	}

	// Check if in GOPATH or GOBIN
	execPath, err := os.Executable()
	if err == nil {
		gopath := os.Getenv("GOPATH")
		gobin := os.Getenv("GOBIN")

		if gopath != "" && strings.Contains(execPath, gopath) {
			return InstallGoInstall
		}
		if gobin != "" && strings.Contains(execPath, gobin) {
			return InstallGoInstall
		}

		// Check default GOPATH
		homeDir, err := os.UserHomeDir()
		if err == nil {
			defaultGoPath := homeDir + "/go/bin"
			if strings.Contains(execPath, defaultGoPath) {
				return InstallGoInstall
			}
		}
	}

	// Default to binary install
	return InstallBinary
}

// GetUpdateInstructions returns update instructions based on install method
func GetUpdateInstructions(method InstallMethod) string {
	switch method {
	case InstallHomebrew:
		return "To update, run:\n  brew upgrade nathan-nicholson/tap/note"
	case InstallGoInstall:
		return "To update, run:\n  go install github.com/nathan-nicholson/note@latest"
	case InstallBinary:
		return "To update, download the latest binary from:\n  https://github.com/nathan-nicholson/note/releases/latest"
	default:
		return "To update, visit:\n  https://github.com/nathan-nicholson/note/releases/latest"
	}
}

// PerformUpdate attempts to update the binary based on install method
func PerformUpdate(method InstallMethod) error {
	switch method {
	case InstallHomebrew:
		return runCommand("brew", "upgrade", "nathan-nicholson/tap/note")
	case InstallGoInstall:
		return runCommand("go", "install", "github.com/nathan-nicholson/note@latest")
	default:
		return fmt.Errorf("automatic update not supported for this installation method")
	}
}

// Helper functions

func isHomebrewInstall() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}

	// Check if in Homebrew Cellar
	if strings.Contains(execPath, "/Cellar/note/") {
		return true
	}

	// Check if symlinked from Homebrew
	cmd := exec.Command("brew", "--prefix", "note")
	output, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) != "" {
		return true
	}

	return false
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// normalizeVersion removes 'v' or 'V' prefix from version strings
func normalizeVersion(v string) string {
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	return v
}

// compareVersions compares two semantic versions
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int

		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &p1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &p2)
		}

		if p1 > p2 {
			return 1
		}
		if p1 < p2 {
			return -1
		}
	}

	return 0
}
