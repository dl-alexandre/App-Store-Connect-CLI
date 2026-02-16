package cmdtest

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.mozilla.org/pkcs7"
	"howett.net/plist"
)

type localProfileItem struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name,omitempty"`
	TeamID    string `json:"teamId,omitempty"`
	BundleID  string `json:"bundleId,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	Path      string `json:"path"`
	Expired   bool   `json:"expired"`
}

func TestProfilesLocal_InstallListCleanExpired(t *testing.T) {
	run := func(args []string) (string, string, error) {
		root := RootCommand("1.2.3")
		root.FlagSet.SetOutput(io.Discard)

		var runErr error
		stdout, stderr := captureOutput(t, func() {
			if err := root.Parse(args); err != nil {
				t.Fatalf("parse error: %v", err)
			}
			runErr = root.Run(context.Background())
		})
		return stdout, stderr, runErr
	}

	installDir := t.TempDir()

	activeUUID := "00000000-0000-0000-0000-000000000001"
	expiredUUID := "00000000-0000-0000-0000-000000000002"

	activeSource := filepath.Join(t.TempDir(), "active.mobileprovision")
	expiredSource := filepath.Join(t.TempDir(), "expired.mobileprovision")

	activeBytes := buildMobileprovision(t, activeUUID, "Active Profile", "TEAM12345", "com.example.app", time.Now().Add(24*time.Hour))
	expiredBytes := buildMobileprovision(t, expiredUUID, "Expired Profile", "TEAM12345", "com.example.app", time.Now().Add(-24*time.Hour))

	if err := os.WriteFile(activeSource, activeBytes, 0o600); err != nil {
		t.Fatalf("WriteFile(activeSource) error: %v", err)
	}
	if err := os.WriteFile(expiredSource, expiredBytes, 0o600); err != nil {
		t.Fatalf("WriteFile(expiredSource) error: %v", err)
	}

	// Install both profiles.
	_, stderr, err := run([]string{"profiles", "local", "install", "--path", activeSource, "--install-dir", installDir, "--output", "json"})
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if err != nil {
		t.Fatalf("run error: %v", err)
	}

	_, stderr, err = run([]string{"profiles", "local", "install", "--path", expiredSource, "--install-dir", installDir, "--output", "json"})
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if err != nil {
		t.Fatalf("run error: %v", err)
	}

	activeInstalled := filepath.Join(installDir, activeUUID+".mobileprovision")
	expiredInstalled := filepath.Join(installDir, expiredUUID+".mobileprovision")

	// List should include both.
	stdout, stderr, err := run([]string{"profiles", "local", "list", "--install-dir", installDir, "--output", "json"})
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if err != nil {
		t.Fatalf("run error: %v", err)
	}

	var items []localProfileItem
	if err := json.Unmarshal([]byte(stdout), &items); err != nil {
		t.Fatalf("decode list JSON: %v (stdout=%q)", err, stdout)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(items))
	}

	// Dry-run clean should plan deletion but not delete.
	_, stderr, err = run([]string{"profiles", "local", "clean", "--install-dir", installDir, "--expired", "--dry-run", "--output", "json"})
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if _, err := os.Stat(expiredInstalled); err != nil {
		t.Fatalf("expected expired profile to still exist after dry-run, stat error: %v", err)
	}

	// Clean without --confirm should be a usage error.
	_, stderr, err = run([]string{"profiles", "local", "clean", "--install-dir", installDir, "--expired", "--output", "json"})
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected flag.ErrHelp, got %v", err)
	}
	if !strings.Contains(stderr, "Error: --confirm is required") {
		t.Fatalf("expected confirm required error, got %q", stderr)
	}

	// Confirmed clean should delete the expired one only.
	_, stderr, err = run([]string{"profiles", "local", "clean", "--install-dir", installDir, "--expired", "--confirm", "--output", "json"})
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if _, err := os.Stat(expiredInstalled); err == nil {
		t.Fatalf("expected expired profile to be deleted")
	}
	if _, err := os.Stat(activeInstalled); err != nil {
		t.Fatalf("expected active profile to remain, stat error: %v", err)
	}
}

func TestProfilesLocalInstall_ByID_DownloadsAndInstalls(t *testing.T) {
	setupAuth(t)

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	uuid := "00000000-0000-0000-0000-0000000000AA"
	content := buildMobileprovision(t, uuid, "Downloaded Profile", "TEAM12345", "com.example.app", time.Now().Add(24*time.Hour))
	b64 := base64.StdEncoding.EncodeToString(content)

	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.appstoreconnect.apple.com" {
			t.Fatalf("unexpected host: %s", req.URL.Host)
		}
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/profiles/p1" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}

		body := `{"data":{"type":"profiles","id":"p1","attributes":{"profileContent":"` + b64 + `"}}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})

	installDir := t.TempDir()

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"profiles", "local", "install", "--id", "p1", "--install-dir", installDir, "--output", "json"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	destPath := filepath.Join(installDir, uuid+".mobileprovision")
	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("ReadFile(destPath) error: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("expected installed profile bytes to match downloaded bytes")
	}
}

func buildMobileprovision(t *testing.T, uuid, name, teamID, bundleID string, expires time.Time) []byte {
	t.Helper()

	now := time.Now().UTC()
	payload := map[string]any{
		"UUID":           uuid,
		"Name":           name,
		"TeamIdentifier": []string{teamID},
		"CreationDate":   now.Add(-1 * time.Hour),
		"ExpirationDate": expires.UTC(),
		"Entitlements": map[string]any{
			"application-identifier":              teamID + "." + bundleID,
			"com.apple.developer.team-identifier": teamID,
		},
	}

	plistBytes, err := plist.Marshal(payload, plist.XMLFormat)
	if err != nil {
		t.Fatalf("plist.Marshal() error: %v", err)
	}

	cert, key := selfSignedCert(t)
	sd, err := pkcs7.NewSignedData(plistBytes)
	if err != nil {
		t.Fatalf("pkcs7.NewSignedData() error: %v", err)
	}
	if err := sd.AddSigner(cert, key, pkcs7.SignerInfoConfig{}); err != nil {
		t.Fatalf("SignedData.AddSigner() error: %v", err)
	}
	out, err := sd.Finish()
	if err != nil {
		t.Fatalf("SignedData.Finish() error: %v", err)
	}
	return out
}

func selfSignedCert(t *testing.T) (*x509.Certificate, crypto.PrivateKey) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey() error: %v", err)
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatalf("rand.Int() error: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: "Test Signer",
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("x509.CreateCertificate() error: %v", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("x509.ParseCertificate() error: %v", err)
	}
	return cert, key
}
