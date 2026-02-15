package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir parent dir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}

func withWorkingDirectory(t *testing.T, path string) {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current working directory: %v", err)
	}
	if err := os.Chdir(path); err != nil {
		t.Fatalf("chdir to %s: %v", path, err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(cwd); chdirErr != nil {
			t.Fatalf("restore working directory: %v", chdirErr)
		}
	})
}

func readJSONEntries(t *testing.T, path string) []wallEntry {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read JSON file %s: %v", path, err)
	}
	var entries []wallEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		t.Fatalf("unmarshal JSON file %s: %v", path, err)
	}
	return entries
}

func TestRunAddsAppAndSyncsReadme(t *testing.T) {
	tmpRepo := t.TempDir()
	withWorkingDirectory(t, tmpRepo)

	writeFile(t, filepath.Join(tmpRepo, "docs", "wall-of-apps.json"), `[
  {
    "app": "CodexMonitor",
    "link": "https://github.com/Dimillian/CodexMonitor",
    "creator": "Dimillian",
    "platform": ["macOS", "iOS"]
  }
]`)

	writeFile(t, filepath.Join(tmpRepo, "README.md"), `# Demo
<!-- WALL-OF-APPS:START -->
Old wall content.
<!-- WALL-OF-APPS:END -->
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := run([]string{
		"--app", "Dandelion",
		"--link", "https://apps.apple.com/us/app/dandelion-write-and-let-go/id6757363901",
		"--creator", "joeycast",
		"--platform", "iOS, macOS",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v (stderr: %s)", err, stderr.String())
	}

	entries := readJSONEntries(t, filepath.Join(tmpRepo, "docs", "wall-of-apps.json"))
	if len(entries) != 2 {
		t.Fatalf("expected 2 JSON entries, got %d", len(entries))
	}
	added := entries[1]
	if added.App != "Dandelion" || added.Creator != "joeycast" {
		t.Fatalf("unexpected added entry: %+v", added)
	}
	if len(added.Platform) != 2 || added.Platform[0] != "iOS" || added.Platform[1] != "macOS" {
		t.Fatalf("unexpected platform values in added entry: %+v", added.Platform)
	}

	readmeBytes, err := os.ReadFile(filepath.Join(tmpRepo, "README.md"))
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	readme := string(readmeBytes)
	expectedRow := "| Dandelion | [Open](https://apps.apple.com/us/app/dandelion-write-and-let-go/id6757363901) | joeycast | iOS, macOS |"
	if !strings.Contains(readme, expectedRow) {
		t.Fatalf("expected generated README row, got:\n%s", readme)
	}
	if !strings.Contains(stdout.String(), "Added app entry in") {
		t.Fatalf("expected add confirmation in stdout, got: %s", stdout.String())
	}
}

func TestRunUpdatesExistingEntryByApp(t *testing.T) {
	tmpRepo := t.TempDir()
	withWorkingDirectory(t, tmpRepo)

	writeFile(t, filepath.Join(tmpRepo, "docs", "wall-of-apps.json"), `[
  {
    "app": "Dandelion",
    "link": "https://apps.apple.com/us/app/dandelion-write-and-let-go/id6757363901",
    "creator": "old-creator",
    "platform": ["iOS"]
  }
]`)

	writeFile(t, filepath.Join(tmpRepo, "README.md"), `# Demo
<!-- WALL-OF-APPS:START -->
Old wall content.
<!-- WALL-OF-APPS:END -->
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := run([]string{
		"--app", "Dandelion",
		"--link", "https://apps.apple.com/us/app/dandelion-write-and-let-go/id6757363901",
		"--creator", "joeycast",
		"--platform", "iOS, macOS",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v (stderr: %s)", err, stderr.String())
	}

	entries := readJSONEntries(t, filepath.Join(tmpRepo, "docs", "wall-of-apps.json"))
	if len(entries) != 1 {
		t.Fatalf("expected single JSON entry after update, got %d", len(entries))
	}
	if entries[0].Creator != "joeycast" {
		t.Fatalf("expected updated creator, got %q", entries[0].Creator)
	}
	if len(entries[0].Platform) != 2 {
		t.Fatalf("expected updated platforms, got %+v", entries[0].Platform)
	}
	if !strings.Contains(stdout.String(), "Updated app entry in") {
		t.Fatalf("expected update confirmation in stdout, got: %s", stdout.String())
	}
}

func TestRunFailsWhenPlatformMissing(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := run([]string{
		"--app", "Dandelion",
		"--link", "https://apps.apple.com/us/app/dandelion-write-and-let-go/id6757363901",
		"--creator", "joeycast",
	}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error when platform is missing")
	}
	if !strings.Contains(err.Error(), "--platform is required") {
		t.Fatalf("expected missing platform error, got %v", err)
	}
}
