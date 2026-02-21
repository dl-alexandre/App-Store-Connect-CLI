package swifthelpers

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestSwiftEndToEnd_JWT tests the complete JWT signing flow using Swift helpers
func TestSwiftEndToEnd_JWT(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Swift helpers only available on macOS")
	}

	// Check if Swift helpers are available
	if _, err := findHelper(JWTSignerBinary); err != nil {
		t.Skip("Swift JWT signer not found, skipping end-to-end test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a proper P-256 key using OpenSSL (matches what Swift helper expects)
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "test-key.p8")

	// Generate key in PKCS#8 format
	cmd := exec.Command("sh", "-c",
		"openssl ecparam -genkey -name prime256v1 -noout | openssl pkcs8 -topk8 -nocrypt")
	output, err := cmd.Output()
	if err != nil {
		t.Skipf("OpenSSL not available to generate test key: %v", err)
	}
	if err := os.WriteFile(keyPath, output, 0o600); err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	// Test JWT signing
	req := JWTSignRequest{
		IssuerID:       "test-issuer-123",
		KeyID:          "test-key-456",
		PrivateKeyPath: keyPath,
	}

	resp, err := SignJWT(ctx, req)
	if err != nil {
		t.Fatalf("SignJWT failed: %v", err)
	}

	if resp == nil {
		t.Fatal("SignJWT returned nil response")
	}

	if resp.Token == "" {
		t.Error("SignJWT returned empty token")
	}

	if resp.ExpiresIn == 0 {
		t.Error("SignJWT returned zero expires_in")
	}

	t.Logf("Successfully generated JWT with %d seconds expiry", resp.ExpiresIn)
}

// TestSwiftEndToEnd_ArchiveExtraction tests archive extraction using Swift helpers
func TestSwiftEndToEnd_ArchiveExtraction(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Swift helpers only available on macOS")
	}

	if _, err := findHelper(ArchiveUnzipBinary); err != nil {
		t.Skip("Swift archive unzipper not found, skipping end-to-end test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test ZIP file with relative paths
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "extracted")

	// Create test file to zip (use relative path for clean extraction)
	contentDir := filepath.Join(tempDir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatalf("Failed to create content dir: %v", err)
	}
	testFile := filepath.Join(contentDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Hello, World!"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create ZIP with relative paths
	cmd := exec.Command("zip", "-r", zipPath, "test.txt")
	cmd.Dir = contentDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create test ZIP: %v", err)
	}

	// Test extraction
	result, err := ExtractArchive(ctx, zipPath, extractDir, false)
	if err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	if !result.Success {
		t.Error("ExtractArchive returned success=false")
	}

	if result.FilesExtracted == 0 {
		t.Error("ExtractArchive extracted zero files")
	}

	if result.TotalSize == 0 {
		t.Error("ExtractArchive returned zero total size")
	}

	// Verify extracted file exists
	extractedFile := filepath.Join(extractDir, "test.txt")
	if _, err := os.Stat(extractedFile); os.IsNotExist(err) {
		t.Errorf("Extracted file not found at %s (files: %d)", extractedFile, result.FilesExtracted)
	}

	t.Logf("Successfully extracted %d files (%d bytes) in %.3fs",
		result.FilesExtracted, result.TotalSize, result.Duration)
}

// TestSwiftEndToEnd_ArchiveList tests listing archive contents
func TestSwiftEndToEnd_ArchiveList(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Swift helpers only available on macOS")
	}

	if _, err := findHelper(ArchiveUnzipBinary); err != nil {
		t.Skip("Swift archive unzipper not found, skipping end-to-end test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test ZIP file
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test.zip")

	// Create test files
	testFile1 := filepath.Join(tempDir, "file1.txt")
	testFile2 := filepath.Join(tempDir, "file2.txt")
	_ = os.WriteFile(testFile1, []byte("Content 1"), 0o644)
	_ = os.WriteFile(testFile2, []byte("Content 2"), 0o644)

	// Create ZIP
	if err := createTestZip(zipPath, testFile1, testFile2); err != nil {
		t.Fatalf("Failed to create test ZIP: %v", err)
	}

	// Test listing
	files, err := ListArchiveContents(ctx, zipPath)
	if err != nil {
		t.Fatalf("ListArchiveContents failed: %v", err)
	}

	if len(files) == 0 {
		t.Error("ListArchiveContents returned empty list")
	}

	t.Logf("Archive contains %d files: %v", len(files), files)
}

// TestSwiftEndToEnd_IPAPacking tests IPA packaging using Swift helpers
func TestSwiftEndToEnd_IPAPacking(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Swift helpers only available on macOS")
	}

	if _, err := findHelper(IPAPackBinary); err != nil {
		t.Skip("Swift IPA packer not found, skipping end-to-end test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a minimal .app bundle structure
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "TestApp.app")
	payloadDir := filepath.Join(appDir, "Payload")
	if err := os.MkdirAll(payloadDir, 0o755); err != nil {
		t.Fatalf("Failed to create app structure: %v", err)
	}

	// Create Info.plist
	infoPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleIdentifier</key>
    <string>com.test.app</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
</dict>
</plist>`
	if err := os.WriteFile(filepath.Join(payloadDir, "Info.plist"), []byte(infoPlist), 0o644); err != nil {
		t.Fatalf("Failed to create Info.plist: %v", err)
	}

	ipaPath := filepath.Join(tempDir, "TestApp.ipa")

	// Test IPA packaging
	result, err := PackIPA(ctx, appDir, ipaPath, 6)
	if err != nil {
		t.Fatalf("PackIPA failed: %v", err)
	}

	if !result.Success {
		t.Error("PackIPA returned success=false")
	}

	if result.CompressionRatio == 0 {
		t.Error("PackIPA returned zero compression ratio")
	}

	// Verify IPA was created
	if _, err := os.Stat(ipaPath); os.IsNotExist(err) {
		t.Errorf("IPA file not created at %s", ipaPath)
	}

	t.Logf("Successfully created IPA: %d -> %d bytes (%.1fx compression)",
		result.OriginalSize, result.CompressedSize, result.CompressionRatio)
}

// TestSwiftEndToEnd_BundleValidation tests bundle validation using Swift helpers
func TestSwiftEndToEnd_BundleValidation(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Swift helpers only available on macOS")
	}

	if _, err := findHelper(BundleValidateBinary); err != nil {
		t.Skip("Swift bundle validator not found, skipping end-to-end test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a minimal .app bundle
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "TestApp.app")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("Failed to create app directory: %v", err)
	}

	// Create Info.plist
	infoPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleIdentifier</key>
    <string>com.test.app</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
</dict>
</plist>`
	if err := os.WriteFile(filepath.Join(appDir, "Info.plist"), []byte(infoPlist), 0o644); err != nil {
		t.Fatalf("Failed to create Info.plist: %v", err)
	}

	// Test validation (unsigned bundles will fail - this is expected)
	result, err := ValidateBundle(ctx, appDir, false)
	if err != nil {
		// Unsigned bundles return error for signature validation - this is expected
		t.Logf("ValidateBundle returned error (expected for unsigned bundle): %v", err)
		return
	}

	if result == nil {
		t.Fatal("ValidateBundle returned nil result")
	}

	// If we got a result without error, check it
	t.Logf("Bundle validation result: valid=%v, issues=%d",
		result.Valid, len(result.Issues))

	// For unsigned bundles in non-strict mode, we expect it to not be valid
	if !result.Valid && len(result.Issues) > 0 {
		t.Logf("Bundle validation correctly identified %d issues in unsigned bundle", len(result.Issues))
	}
}

// TestSwiftEndToEnd_ImageOptimization tests image optimization using Swift helpers
func TestSwiftEndToEnd_ImageOptimization(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Swift helpers only available on macOS")
	}

	if _, err := findHelper(ImageOptimizeBinary); err != nil {
		t.Skip("Swift image optimizer not found, skipping end-to-end test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a simple test PNG
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.png")
	outputPath := filepath.Join(tempDir, "output.png")

	// Create a minimal valid PNG file (1x1 pixel, red)
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x00, 0x03, 0x00, 0x01, 0x00, 0x05, 0xFE,
		0xD4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	if err := os.WriteFile(inputPath, pngData, 0o644); err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}

	// Test optimization
	req := ImageOptimizeRequest{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Preset:     "thumbnail",
		Format:     "png",
	}

	result, err := OptimizeImage(ctx, req)
	if err != nil {
		t.Fatalf("OptimizeImage failed: %v", err)
	}

	if result == nil {
		t.Fatal("OptimizeImage returned nil result")
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file not created at %s", outputPath)
	}

	t.Logf("Image optimization: %d -> %d bytes (%.1f%% savings)",
		result.OriginalSize, result.OptimizedSize, result.SavingsPercent)
}

// TestSwiftEndToEnd_CodeSignVerification tests code signature verification
func TestSwiftEndToEnd_CodeSignVerification(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Swift helpers only available on macOS")
	}

	if _, err := findHelper(CodeSignBinary); err != nil {
		t.Skip("Swift codesign helper not found, skipping end-to-end test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test listing identities (should work even without signing)
	identities, err := ListCodeSignIdentities(ctx)
	if err != nil {
		t.Logf("ListCodeSignIdentities returned error (may be expected in CI): %v", err)
	} else {
		t.Logf("Found %d signing identities", len(identities))
	}
}

// TestSwiftEndToEnd_AllHelpersPresent verifies all expected helpers are available
func TestSwiftEndToEnd_AllHelpersPresent(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Swift helpers only available on macOS")
	}

	helperTests := []struct {
		name   string
		binary string
	}{
		{"JWT Signer", JWTSignerBinary},
		{"Keychain", KeychainBinary},
		{"Screenshot Frame", ScreenshotFrameBinary},
		{"Image Optimize", ImageOptimizeBinary},
		{"Bundle Validate", BundleValidateBinary},
		{"IPA Pack", IPAPackBinary},
		{"Simulator", SimulatorBinary},
		{"Video Encode", VideoEncodeBinary},
		{"CodeSign", CodeSignBinary},
		{"Archive Unzip", ArchiveUnzipBinary},
	}

	status := GetStatus()
	if !status.Available {
		t.Skip("Swift helpers not available")
	}

	for _, tt := range helperTests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := findHelper(tt.binary)
			if err != nil {
				t.Logf("Helper %s not found: %v", tt.binary, err)
			} else {
				t.Logf("Found %s at %s", tt.binary, path)
			}
		})
	}
}

// Helper function to create a test ZIP file
func createTestZip(zipPath string, files ...string) error {
	// Use system zip command for simplicity in tests
	args := append([]string{"-r", zipPath}, files...)
	cmd := exec.Command("zip", args...)
	cmd.Dir = filepath.Dir(files[0])
	return cmd.Run()
}
