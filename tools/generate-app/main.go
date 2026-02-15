package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/wallgen"
)

type wallEntry struct {
	App      string   `json:"app"`
	Link     string   `json:"link"`
	Creator  string   `json:"creator"`
	Platform []string `json:"platform"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer, stderr io.Writer) error {
	fs := flag.NewFlagSet("generate-app", flag.ContinueOnError)
	fs.SetOutput(stderr)

	app := fs.String("app", "", "App name")
	link := fs.String("link", "", "App URL")
	creator := fs.String("creator", "", "Creator name or handle")
	platformCSV := fs.String("platform", "", "Comma-separated platform labels (for example: iOS,macOS)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(fs.Args(), ", "))
	}

	entry, err := buildEntry(*app, *link, *creator, *platformCSV)
	if err != nil {
		return err
	}

	repoRoot, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	jsonPath := sourcePath(repoRoot)

	entries, err := readEntries(jsonPath)
	if err != nil {
		return err
	}
	entries, action := upsertByApp(entries, entry)

	if err := writeEntries(jsonPath, entries); err != nil {
		return err
	}
	result, err := wallgen.Generate(repoRoot)
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "%s app entry in %s\n", action, jsonPath)
	fmt.Fprintf(stdout, "Synced snippet markers in %s\n", result.ReadmePath)
	return nil
}

func buildEntry(app, link, creator, platformCSV string) (wallEntry, error) {
	app = strings.TrimSpace(app)
	link = strings.TrimSpace(link)
	creator = strings.TrimSpace(creator)

	if app == "" {
		return wallEntry{}, fmt.Errorf("--app is required")
	}
	if link == "" {
		return wallEntry{}, fmt.Errorf("--link is required")
	}
	if creator == "" {
		return wallEntry{}, fmt.Errorf("--creator is required")
	}
	if err := validateHTTPURL(link); err != nil {
		return wallEntry{}, fmt.Errorf("--link must be a valid http/https URL")
	}

	platforms := splitPlatforms(platformCSV)
	if len(platforms) == 0 {
		return wallEntry{}, fmt.Errorf("--platform is required")
	}

	return wallEntry{
		App:      app,
		Link:     link,
		Creator:  creator,
		Platform: platforms,
	}, nil
}

func sourcePath(repoRoot string) string {
	return filepath.Join(repoRoot, "docs", "wall-of-apps.json")
}

func readEntries(path string) ([]wallEntry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("missing source file: %s", path)
	}
	if strings.TrimSpace(string(raw)) == "" {
		return nil, fmt.Errorf("source file is empty: %s", path)
	}

	var entries []wallEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("source file has no entries: %s", path)
	}
	return entries, nil
}

func upsertByApp(entries []wallEntry, candidate wallEntry) ([]wallEntry, string) {
	for i := range entries {
		if strings.EqualFold(strings.TrimSpace(entries[i].App), candidate.App) {
			entries[i] = candidate
			return entries, "Updated"
		}
	}
	return append(entries, candidate), "Added"
}

func writeEntries(path string, entries []wallEntry) error {
	content, err := renderEntries(entries)
	if err != nil {
		return fmt.Errorf("marshal source JSON: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write source JSON: %w", err)
	}
	return nil
}

func renderEntries(entries []wallEntry) (string, error) {
	var builder strings.Builder
	builder.WriteString("[\n")
	for i, entry := range entries {
		app, err := quoteJSON(entry.App)
		if err != nil {
			return "", err
		}
		link, err := quoteJSON(entry.Link)
		if err != nil {
			return "", err
		}
		creator, err := quoteJSON(entry.Creator)
		if err != nil {
			return "", err
		}

		builder.WriteString("  {\n")
		builder.WriteString("    \"app\": ")
		builder.WriteString(app)
		builder.WriteString(",\n")
		builder.WriteString("    \"link\": ")
		builder.WriteString(link)
		builder.WriteString(",\n")
		builder.WriteString("    \"creator\": ")
		builder.WriteString(creator)
		builder.WriteString(",\n")
		builder.WriteString("    \"platform\": [")
		for platformIndex, platform := range entry.Platform {
			quotedPlatform, err := quoteJSON(platform)
			if err != nil {
				return "", err
			}
			if platformIndex > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(quotedPlatform)
		}
		builder.WriteString("]\n")
		builder.WriteString("  }")
		if i < len(entries)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}
	builder.WriteString("]\n")
	return builder.String(), nil
}

func quoteJSON(value string) (string, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return "", err
	}
	return strings.TrimSuffix(buffer.String(), "\n"), nil
}

func splitPlatforms(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	platforms := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" || containsFold(platforms, trimmed) {
			continue
		}
		platforms = append(platforms, trimmed)
	}
	return platforms
}

func containsFold(values []string, needle string) bool {
	for _, value := range values {
		if strings.EqualFold(value, needle) {
			return true
		}
	}
	return false
}

func validateHTTPURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("invalid URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme")
	}
	return nil
}
