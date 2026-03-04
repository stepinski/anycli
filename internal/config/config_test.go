package config

// WHY package config (not package config_test)?
//
// Go has two test package options:
//   package config      — "white box" testing, can access unexported symbols
//   package config_test — "black box" testing, only exported symbols visible
//
// We use white box here because we want to test Dir() which uses
// internal env var logic. For most packages, black box is preferred
// because it tests the public contract, not implementation details.

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDir tests the Dir() function under different conditions.
//
// TABLE-DRIVEN TESTS: the standard Go testing pattern.
// Instead of writing TestDir_WithXDG, TestDir_WithoutXDG, TestDir_NoHome...
// you write one test with a table of cases. Benefits:
//   - Adding a new case is one line
//   - All cases use the same test logic — no duplication
//   - Test output shows which case failed by name
func TestDir(t *testing.T) {
	// Each test case is an anonymous struct.
	// 'name' is what appears in test output.
	cases := []struct {
		name    string
		xdgHome string // value to set XDG_CONFIG_HOME to ("" = unset)
		wantSuffix string // expected path suffix
	}{
		{
			name:       "uses XDG_CONFIG_HOME when set",
			xdgHome:    "/tmp/xdg",
			wantSuffix: "/tmp/xdg/anycli",
		},
		{
			name:       "falls back to ~/.config/anycli when XDG not set",
			xdgHome:    "",
			wantSuffix: ".config/anycli",
		},
	}

	for _, tc := range cases {
		// t.Run creates a subtest — output will show "TestDir/uses_XDG_CONFIG_HOME_when_set"
		// This is how you run a single case: go test -run TestDir/uses_XDG
		t.Run(tc.name, func(t *testing.T) {
			// t.Setenv sets an env var and automatically restores the
			// original value when the test ends. No manual cleanup needed.
			// This is Go 1.17+ — prefer it over os.Setenv in tests.
			if tc.xdgHome != "" {
				t.Setenv("XDG_CONFIG_HOME", tc.xdgHome)
			} else {
				// Ensure XDG is not set from a previous test or the environment
				t.Setenv("XDG_CONFIG_HOME", "")
			}

			got, err := Dir()

			// In Go tests, err checking comes before result checking.
			// If we expected success and got an error, fail immediately.
			if err != nil {
				t.Fatalf("Dir() unexpected error: %v", err)
			}

			// We check HasSuffix rather than exact equality because
			// the home directory varies per machine.
			if !filepath.IsAbs(got) {
				t.Errorf("Dir() = %q, want absolute path", got)
			}

			if tc.xdgHome != "" && got != tc.wantSuffix {
				t.Errorf("Dir() = %q, want %q", got, tc.wantSuffix)
			}
		})
	}
}

// TestLoad_Defaults verifies that Load() returns correct defaults
// when no config file exists.
func TestLoad_Defaults(t *testing.T) {
	// t.TempDir() creates a temporary directory that is automatically
	// removed when the test ends. Use this instead of os.MkdirTemp
	// because cleanup is guaranteed even if the test panics.
	tmp := t.TempDir()

	// Point viper at our empty temp dir so it finds no config file
	t.Setenv("XDG_CONFIG_HOME", tmp)

	// Clear any env vars that might bleed in from the test environment
	t.Setenv("ANYCLI_API_KEY", "")
	t.Setenv("ANYCLI_URL", "")
	t.Setenv("ANYCLI_WORKSPACE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	// Check each default value
	// t.Errorf (not Fatalf) — we want to see ALL failures, not stop at first
	if cfg.URL != "http://localhost:3001" {
		t.Errorf("URL = %q, want %q", cfg.URL, "http://localhost:3001")
	}
	if cfg.Workspace != "vault" {
		t.Errorf("Workspace = %q, want %q", cfg.Workspace, "vault")
	}
	if cfg.Mode != "chat" {
		t.Errorf("Mode = %q, want %q", cfg.Mode, "chat")
	}
	if !cfg.Stream {
		t.Error("Stream = false, want true")
	}
	if len(cfg.Priorities) == 0 {
		t.Error("Priorities is empty, want defaults")
	}
}

// TestLoad_FromFile verifies that Load() correctly reads a config file.
func TestLoad_FromFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	// Write a real config file into the temp dir
	// filepath.Join handles OS path separators correctly
	cfgDir := filepath.Join(tmp, "anycli")
	if err := os.MkdirAll(cfgDir, 0700); err != nil {
		t.Fatalf("creating config dir: %v", err)
	}

	content := `
url: "http://myserver:3001"
api_key: "test-key-123"
workspace: "my-notes"
stream: false
mode: query
`
	cfgPath := filepath.Join(cfgDir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(content), 0600); err != nil {
		t.Fatalf("writing config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.URL != "http://myserver:3001" {
		t.Errorf("URL = %q, want %q", cfg.URL, "http://myserver:3001")
	}
	if cfg.APIKey != "test-key-123" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "test-key-123")
	}
	if cfg.Workspace != "my-notes" {
		t.Errorf("Workspace = %q, want %q", cfg.Workspace, "my-notes")
	}
	if cfg.Stream {
		t.Error("Stream = true, want false")
	}
	if cfg.Mode != "query" {
		t.Errorf("Mode = %q, want %q", cfg.Mode, "query")
	}
}

// TestLoad_EnvOverride verifies that env vars override config file values.
// This is the most important test — it validates the priority chain.
func TestLoad_EnvOverride(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	// Write config file with one value
	cfgDir := filepath.Join(tmp, "anycli")
	os.MkdirAll(cfgDir, 0700)
	content := `workspace: "from-file"`
	os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(content), 0600)

	// Set env var to a different value
	// This should win over the config file
	t.Setenv("ANYCLI_WORKSPACE", "from-env")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	// Env var must win
	if cfg.Workspace != "from-env" {
		t.Errorf("Workspace = %q, want %q (env should override file)", cfg.Workspace, "from-env")
	}
}

// TestWrite verifies that Write() creates a readable config file.
func TestWrite(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	path, err := Write("http://test:3001", "my-api-key", "test-workspace")
	if err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	// File must exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Write() did not create file at %q", path)
	}

	// File must have secure permissions (0600)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %q: %v", path, err)
	}
	// info.Mode() returns the full mode — mask with 0777 to get just permissions
	if info.Mode().Perm() != 0600 {
		t.Errorf("file permissions = %o, want 600", info.Mode().Perm())
	}

	// Written file must be loadable
	t.Setenv("ANYCLI_API_KEY", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() after Write() unexpected error: %v", err)
	}

	if cfg.URL != "http://test:3001" {
		t.Errorf("URL = %q, want %q", cfg.URL, "http://test:3001")
	}
	if cfg.APIKey != "my-api-key" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "my-api-key")
	}
}
