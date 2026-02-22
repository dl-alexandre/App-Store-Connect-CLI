// Package syscli provides Go wrappers around system CLIs for macOS
// operations that require native tooling (codesign, simctl, security).
//
// All operations call system binaries directly from Go — no Swift helper
// subprocesses are involved. JWT signing and archive/zip operations are
// handled entirely in-process by the Go standard library and golang-jwt.
package syscli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// IsDarwin reports whether the current platform is macOS.
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// runCmd executes a command and returns stdout. On error it includes stderr
// in the returned error message.
func runCmd(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s failed: %w (stderr: %s)", name, err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

// --- Simulator (xcrun simctl) ---

// SimulatorDevice represents an iOS simulator.
type SimulatorDevice struct {
	UDID        string `json:"udid"`
	Name        string `json:"name"`
	DeviceType  string `json:"deviceTypeIdentifier"`
	Runtime     string `json:"runtime"`
	State       string `json:"state"`
	IsAvailable bool   `json:"isAvailable"`
}

// ListSimulators returns available iOS simulators via xcrun simctl.
func ListSimulators(ctx context.Context, bootedOnly bool) ([]SimulatorDevice, error) {
	if !IsDarwin() {
		return nil, fmt.Errorf("simulators not available on %s", runtime.GOOS)
	}

	out, err := runCmd(ctx, "xcrun", "simctl", "list", "devices", "-j")
	if err != nil {
		return nil, err
	}

	var result struct {
		Devices map[string][]SimulatorDevice `json:"devices"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("failed to parse simctl output: %w", err)
	}

	var devices []SimulatorDevice
	for runtimeID, devs := range result.Devices {
		for _, d := range devs {
			if bootedOnly && d.State != "Booted" {
				continue
			}
			d.Runtime = runtimeID
			devices = append(devices, d)
		}
	}
	return devices, nil
}

// SimulatorScreenshot captures a screenshot from a simulator.
func SimulatorScreenshot(ctx context.Context, deviceUDID, outputPath string) error {
	if !IsDarwin() {
		return fmt.Errorf("simulators not available on %s", runtime.GOOS)
	}
	_, err := runCmd(ctx, "xcrun", "simctl", "io", deviceUDID, "screenshot", outputPath)
	return err
}

// SimulatorInstall installs an app on a simulator.
func SimulatorInstall(ctx context.Context, deviceUDID, appPath string) error {
	if !IsDarwin() {
		return fmt.Errorf("simulators not available on %s", runtime.GOOS)
	}
	_, err := runCmd(ctx, "xcrun", "simctl", "install", deviceUDID, appPath)
	return err
}

// SimulatorLaunch launches an app on a simulator.
func SimulatorLaunch(ctx context.Context, deviceUDID, bundleID string) error {
	if !IsDarwin() {
		return fmt.Errorf("simulators not available on %s", runtime.GOOS)
	}
	_, err := runCmd(ctx, "xcrun", "simctl", "launch", deviceUDID, bundleID)
	return err
}

// --- Code Signing (/usr/bin/codesign) ---

// CodeSignResult is returned after code signing.
type CodeSignResult struct {
	Success bool   `json:"success"`
	Path    string `json:"path"`
}

// CodeSign signs an app bundle using /usr/bin/codesign.
func CodeSign(ctx context.Context, path, identity, entitlements string, force bool) (*CodeSignResult, error) {
	if !IsDarwin() {
		return nil, fmt.Errorf("codesign not available on %s", runtime.GOOS)
	}

	args := []string{"--sign"}
	if identity != "" {
		args = append(args, identity)
	} else {
		args = append(args, "-")
	}
	if entitlements != "" {
		args = append(args, "--entitlements", entitlements)
	}
	if force {
		args = append(args, "--force")
	}
	args = append(args, path)

	_, err := runCmd(ctx, "/usr/bin/codesign", args...)
	if err != nil {
		return nil, err
	}
	return &CodeSignResult{Success: true, Path: path}, nil
}

// CodeSignVerifyResult is returned after signature verification.
type CodeSignVerifyResult struct {
	Valid          bool   `json:"valid"`
	Path           string `json:"path"`
	Authority      string `json:"authority,omitempty"`
	Identifier     string `json:"identifier,omitempty"`
	TeamIdentifier string `json:"teamIdentifier,omitempty"`
}

// CodeSignVerify verifies code signature using /usr/bin/codesign.
func CodeSignVerify(ctx context.Context, path string) (*CodeSignVerifyResult, error) {
	if !IsDarwin() {
		return nil, fmt.Errorf("codesign not available on %s", runtime.GOOS)
	}

	out, err := runCmd(ctx, "/usr/bin/codesign", "--verify", "--verbose=2", path)
	if err != nil {
		return &CodeSignVerifyResult{Valid: false, Path: path}, nil
	}

	result := &CodeSignVerifyResult{Valid: true, Path: path}

	// Parse verbose output for authority info
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Authority=") {
			result.Authority = strings.TrimPrefix(line, "Authority=")
		}
		if strings.HasPrefix(line, "Identifier=") {
			result.Identifier = strings.TrimPrefix(line, "Identifier=")
		}
		if strings.HasPrefix(line, "TeamIdentifier=") {
			result.TeamIdentifier = strings.TrimPrefix(line, "TeamIdentifier=")
		}
	}
	return result, nil
}

// --- Security (/usr/bin/security) ---

// FindIdentity lists available code signing identities via /usr/bin/security.
func FindIdentity(ctx context.Context, policy string) ([]string, error) {
	if !IsDarwin() {
		return nil, fmt.Errorf("security not available on %s", runtime.GOOS)
	}

	args := []string{"find-identity", "-v"}
	if policy != "" {
		args = append(args, "-p", policy)
	}

	out, err := runCmd(ctx, "/usr/bin/security", args...)
	if err != nil {
		return nil, err
	}

	var identities []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Policy:") && !strings.Contains(line, "valid identities found") {
			identities = append(identities, line)
		}
	}
	return identities, nil
}
