package update

import (
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "version with v prefix",
			input:    "v1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "version without v prefix",
			input:    "1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "version with V prefix uppercase",
			input:    "V2.0.0",
			expected: "2.0.0",
		},
		{
			name:     "empty version",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeVersion(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeVersion(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{
			name:     "equal versions",
			v1:       "1.2.3",
			v2:       "1.2.3",
			expected: 0,
		},
		{
			name:     "v1 greater - major",
			v1:       "2.0.0",
			v2:       "1.9.9",
			expected: 1,
		},
		{
			name:     "v1 greater - minor",
			v1:       "1.5.0",
			v2:       "1.4.9",
			expected: 1,
		},
		{
			name:     "v1 greater - patch",
			v1:       "1.2.4",
			v2:       "1.2.3",
			expected: 1,
		},
		{
			name:     "v2 greater - major",
			v1:       "1.9.9",
			v2:       "2.0.0",
			expected: -1,
		},
		{
			name:     "v2 greater - minor",
			v1:       "1.4.9",
			v2:       "1.5.0",
			expected: -1,
		},
		{
			name:     "v2 greater - patch",
			v1:       "1.2.3",
			v2:       "1.2.4",
			expected: -1,
		},
		{
			name:     "different lengths - v1 longer",
			v1:       "1.2.3.4",
			v2:       "1.2.3",
			expected: 1,
		},
		{
			name:     "different lengths - v2 longer",
			v1:       "1.2.3",
			v2:       "1.2.3.4",
			expected: -1,
		},
		{
			name:     "major version difference",
			v1:       "2.0.0",
			v2:       "1.99.99",
			expected: 1,
		},
		{
			name:     "zero versions",
			v1:       "0.0.0",
			v2:       "0.0.0",
			expected: 0,
		},
		{
			name:     "pre-release to release",
			v1:       "1.0.0",
			v2:       "0.9.9",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestDetectInstallMethod(t *testing.T) {
	// This test is environment-dependent, so we just verify it returns a valid method
	method := DetectInstallMethod()

	validMethods := []InstallMethod{
		InstallUnknown,
		InstallHomebrew,
		InstallGoInstall,
		InstallBinary,
	}

	found := false
	for _, valid := range validMethods {
		if method == valid {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("DetectInstallMethod() returned invalid method: %v", method)
	}
}

func TestGetUpdateInstructions(t *testing.T) {
	tests := []struct {
		name     string
		method   InstallMethod
		expected string
	}{
		{
			name:     "homebrew",
			method:   InstallHomebrew,
			expected: "brew upgrade",
		},
		{
			name:     "go install",
			method:   InstallGoInstall,
			expected: "go install",
		},
		{
			name:     "binary",
			method:   InstallBinary,
			expected: "download",
		},
		{
			name:     "unknown",
			method:   InstallUnknown,
			expected: "visit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUpdateInstructions(tt.method)
			if result == "" {
				t.Error("GetUpdateInstructions() returned empty string")
			}
			// Check that the expected keyword is in the instructions
			if !containsIgnoreCase(result, tt.expected) {
				t.Errorf("GetUpdateInstructions(%v) = %q, expected to contain %q", tt.method, result, tt.expected)
			}
		})
	}
}

func TestUpdateInfo(t *testing.T) {
	info := &UpdateInfo{
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.1.0",
		UpdateURL:      "https://example.com/release",
		IsNewer:        true,
		ReleaseNotes:   "Test release notes",
	}

	if info.CurrentVersion != "1.0.0" {
		t.Errorf("CurrentVersion = %q, want %q", info.CurrentVersion, "1.0.0")
	}

	if info.LatestVersion != "1.1.0" {
		t.Errorf("LatestVersion = %q, want %q", info.LatestVersion, "1.1.0")
	}

	if !info.IsNewer {
		t.Error("IsNewer = false, want true")
	}
}

// Helper function
func containsIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return contains(s, substr)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	return len(substr) <= len(s) && (substr == "" || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
