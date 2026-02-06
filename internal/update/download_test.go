package update

import "testing"

func TestDisplayOS(t *testing.T) {
	tests := []struct {
		name string
		goos string
		want string
	}{
		{name: "darwin maps to macOS", goos: "darwin", want: "macOS"},
		{name: "linux passes through", goos: "linux", want: "linux"},
		{name: "windows passes through", goos: "windows", want: "windows"},
		{name: "empty passes through", goos: "", want: ""},
		{name: "freebsd passes through", goos: "freebsd", want: "freebsd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayOS(tt.goos)
			if got != tt.want {
				t.Errorf("displayOS(%q) = %q, want %q", tt.goos, got, tt.want)
			}
		})
	}
}

func TestAssetName(t *testing.T) {
	tests := []struct {
		name       string
		binaryName string
		osName     string
		arch       string
		version    string
		want       string
	}{
		// ── Supported release platforms ──────────────────────────────
		{
			name:       "macOS Intel",
			binaryName: "asc", osName: "darwin", arch: "amd64", version: "0.25.0",
			want: "asc_0.25.0_macOS_amd64",
		},
		{
			name:       "macOS Apple Silicon",
			binaryName: "asc", osName: "darwin", arch: "arm64", version: "0.25.0",
			want: "asc_0.25.0_macOS_arm64",
		},
		{
			name:       "Linux amd64",
			binaryName: "asc", osName: "linux", arch: "amd64", version: "0.25.0",
			want: "asc_0.25.0_linux_amd64",
		},
		{
			name:       "Linux arm64",
			binaryName: "asc", osName: "linux", arch: "arm64", version: "0.25.0",
			want: "asc_0.25.0_linux_arm64",
		},
		{
			name:       "Windows gets .exe extension",
			binaryName: "asc", osName: "windows", arch: "amd64", version: "0.25.0",
			want: "asc_0.25.0_windows_amd64.exe",
		},

		// ── Field order: binary_version_os_arch ─────────────────────
		{
			name:       "field order is binary_version_os_arch",
			binaryName: "tool", osName: "darwin", arch: "arm64", version: "1.2.3",
			want: "tool_1.2.3_macOS_arm64",
		},

		// ── Version variations ──────────────────────────────────────
		{
			name:       "version with patch",
			binaryName: "asc", osName: "darwin", arch: "arm64", version: "0.24.2",
			want: "asc_0.24.2_macOS_arm64",
		},
		{
			name:       "dev version for CI builds",
			binaryName: "asc", osName: "linux", arch: "amd64", version: "dev",
			want: "asc_dev_linux_amd64",
		},

		// ── Empty field edge cases (all must return "") ─────────────
		{
			name:       "empty binary name",
			binaryName: "", osName: "darwin", arch: "arm64", version: "1.0.0",
			want: "",
		},
		{
			name:       "empty OS",
			binaryName: "asc", osName: "", arch: "arm64", version: "1.0.0",
			want: "",
		},
		{
			name:       "empty arch",
			binaryName: "asc", osName: "darwin", arch: "", version: "1.0.0",
			want: "",
		},
		{
			name:       "empty version",
			binaryName: "asc", osName: "darwin", arch: "arm64", version: "",
			want: "",
		},
		{
			name:       "all empty",
			binaryName: "", osName: "", arch: "", version: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assetName(tt.binaryName, tt.osName, tt.arch, tt.version)
			if got != tt.want {
				t.Errorf("assetName(%q, %q, %q, %q) = %q, want %q",
					tt.binaryName, tt.osName, tt.arch, tt.version, got, tt.want)
			}
		})
	}
}

// TestAssetNameMatchesReleaseWorkflow verifies that the Go function produces
// the same filenames as the release.yml shell commands. If this test fails,
// the self-updater and the release workflow are out of sync.
func TestAssetNameMatchesReleaseWorkflow(t *testing.T) {
	// These are the exact outputs produced by the release workflow for
	// a tag "0.25.0": the shell template asc_${VERSION}_macOS_arm64 etc.
	// If the naming convention changes, update BOTH this test AND release.yml.
	releaseAssets := map[string]string{
		"darwin/amd64":  "asc_0.25.0_macOS_amd64",
		"darwin/arm64":  "asc_0.25.0_macOS_arm64",
		"linux/amd64":   "asc_0.25.0_linux_amd64",
		"linux/arm64":   "asc_0.25.0_linux_arm64",
		"windows/amd64": "asc_0.25.0_windows_amd64.exe",
	}

	for platform, wantAsset := range releaseAssets {
		t.Run(platform, func(t *testing.T) {
			var osName, arch string
			for i, c := range platform {
				if c == '/' {
					osName = platform[:i]
					arch = platform[i+1:]
					break
				}
			}
			got := assetName("asc", osName, arch, "0.25.0")
			if got != wantAsset {
				t.Errorf("assetName(\"asc\", %q, %q, \"0.25.0\") = %q, want %q (must match release.yml)",
					osName, arch, got, wantAsset)
			}
		})
	}
}

func TestParseChecksum(t *testing.T) {
	tests := []struct {
		name  string
		data  string
		asset string
		want  string
	}{
		{
			name:  "finds checksum for asset",
			data:  "abc123  asc_0.25.0_macOS_arm64\ndef456  asc_0.25.0_linux_amd64\n",
			asset: "asc_0.25.0_macOS_arm64",
			want:  "abc123",
		},
		{
			name:  "finds second asset",
			data:  "abc123  asc_0.25.0_macOS_arm64\ndef456  asc_0.25.0_linux_amd64\n",
			asset: "asc_0.25.0_linux_amd64",
			want:  "def456",
		},
		{
			name:  "asset not found returns empty",
			data:  "abc123  asc_0.25.0_macOS_arm64\n",
			asset: "asc_0.25.0_linux_amd64",
			want:  "",
		},
		{
			name:  "empty data returns empty",
			data:  "",
			asset: "asc_0.25.0_macOS_arm64",
			want:  "",
		},
		{
			name:  "handles shasum double-space format",
			data:  "abc123def456789012345678901234567890123456789012345678901234  asc_0.25.0_macOS_arm64\n",
			asset: "asc_0.25.0_macOS_arm64",
			want:  "abc123def456789012345678901234567890123456789012345678901234",
		},
		{
			name:  "ignores lines with insufficient fields",
			data:  "lonely\nabc123  asc_0.25.0_macOS_arm64\n",
			asset: "asc_0.25.0_macOS_arm64",
			want:  "abc123",
		},
		{
			name:  "handles trailing newline",
			data:  "abc123  asc_0.25.0_macOS_arm64\n\n",
			asset: "asc_0.25.0_macOS_arm64",
			want:  "abc123",
		},
		{
			name:  "handles no trailing newline",
			data:  "abc123  asc_0.25.0_macOS_arm64",
			asset: "asc_0.25.0_macOS_arm64",
			want:  "abc123",
		},
		{
			name:  "windows asset with .exe",
			data:  "abc123  asc_0.25.0_windows_amd64.exe\n",
			asset: "asc_0.25.0_windows_amd64.exe",
			want:  "abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseChecksum(tt.data, tt.asset)
			if got != tt.want {
				t.Errorf("parseChecksum(..., %q) = %q, want %q", tt.asset, got, tt.want)
			}
		})
	}
}
