package swifthelpers

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// BenchmarkJWTSigning compares Go (golang-jwt) vs Swift (CryptoKit) JWT signing performance
func BenchmarkJWTSigning(b *testing.B) {
	if runtime.GOOS != "darwin" {
		b.Skip("Swift helpers only available on macOS")
	}

	// Check if Swift helper is available
	_, swiftAvailable := findHelper(JWTSignerBinary)

	// Generate test key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.Fatalf("Failed to generate key: %v", err)
	}

	// Export to temp file for Swift
	tempDir := b.TempDir()
	keyPath := filepath.Join(tempDir, "bench-key.p8")

	privKeyBytes, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyBytes}
	keyFile, _ := os.Create(keyPath)
	_ = pem.Encode(keyFile, block)
	keyFile.Close()

	ctx := context.Background()

	// Benchmark Go implementation
	b.Run("Go_golang-jwt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := generateJWTGo("test-key", "test-issuer", privateKey)
			if err != nil {
				b.Fatalf("Go JWT generation failed: %v", err)
			}
		}
	})

	// Benchmark Swift implementation (if available)
	if swiftAvailable == nil {
		b.Run("Swift_CryptoKit", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := SignJWT(ctx, JWTSignRequest{
					IssuerID:       "test-issuer",
					KeyID:          "test-key",
					PrivateKeyPath: keyPath,
				})
				if err != nil {
					b.Fatalf("Swift JWT generation failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkArchiveExtraction compares Go vs Swift archive extraction
func BenchmarkArchiveExtraction(b *testing.B) {
	if runtime.GOOS != "darwin" {
		b.Skip("Swift helpers only available on macOS")
	}

	_, swiftAvailable := findHelper(ArchiveUnzipBinary)

	// Create test ZIP with varying sizes
	tempDir := b.TempDir()
	sizes := []int{1, 10, 100} // MB

	for _, sizeMB := range sizes {
		// Create test file of specified size
		testFile := filepath.Join(tempDir, fmt.Sprintf("test-%dMB.bin", sizeMB))
		data := make([]byte, sizeMB*1024*1024)
		_ = os.WriteFile(testFile, data, 0o644)

		// Create ZIP
		zipPath := filepath.Join(tempDir, fmt.Sprintf("test-%dMB.zip", sizeMB))
		cmd := exec.Command("zip", "-r", zipPath, testFile)
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			b.Fatalf("Failed to create test ZIP: %v", err)
		}

		extractDir := filepath.Join(tempDir, fmt.Sprintf("extract-%dMB", sizeMB))
		_ = os.MkdirAll(extractDir, 0o755)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Benchmark Swift extraction
		if swiftAvailable == nil {
			b.Run(fmt.Sprintf("Swift_%dMB", sizeMB), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					// Clean extract dir
					_ = os.RemoveAll(extractDir)
					_ = os.MkdirAll(extractDir, 0o755)

					_, err := ExtractArchive(ctx, zipPath, extractDir, false)
					if err != nil {
						b.Fatalf("Swift extraction failed: %v", err)
					}
				}
			})
		}

		// Benchmark Go extraction (using archive/zip)
		b.Run(fmt.Sprintf("Go_%dMB", sizeMB), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Clean extract dir
				_ = os.RemoveAll(extractDir)
				_ = os.MkdirAll(extractDir, 0o755)

				err := extractZipGo(zipPath, extractDir)
				if err != nil {
					b.Fatalf("Go extraction failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkImageOptimization compares Go vs Swift image optimization
func BenchmarkImageOptimization(b *testing.B) {
	if runtime.GOOS != "darwin" {
		b.Skip("Swift helpers only available on macOS")
	}

	_, swiftAvailable := findHelper(ImageOptimizeBinary)

	// Create test PNG images of different sizes
	tempDir := b.TempDir()
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"Small", 100, 100},
		{"Medium", 1000, 1000},
		{"Large", 3000, 3000},
	}

	for _, size := range sizes {
		inputPath := filepath.Join(tempDir, fmt.Sprintf("%s.png", size.name))
		outputPath := filepath.Join(tempDir, fmt.Sprintf("%s-optimized.png", size.name))

		// Create test PNG
		if err := createTestPNG(inputPath, size.width, size.height); err != nil {
			b.Fatalf("Failed to create test PNG: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Benchmark Swift optimization
		if swiftAvailable == nil {
			b.Run(fmt.Sprintf("Swift_%s", size.name), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					os.Remove(outputPath)
					_, err := OptimizeImage(ctx, ImageOptimizeRequest{
						InputPath:  inputPath,
						OutputPath: outputPath,
						Preset:     "store",
						Format:     "png",
					})
					if err != nil {
						b.Fatalf("Swift optimization failed: %v", err)
					}
				}
			})
		}

		// Benchmark Go optimization (basic resize)
		b.Run(fmt.Sprintf("Go_%s", size.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = os.Remove(outputPath)
				// Simple file copy as baseline (Go doesn't have native image optimization)
				data, _ := os.ReadFile(inputPath)
				_ = os.WriteFile(outputPath, data, 0o644)
			}
		})
	}
}

// BenchmarkIPAPackaging compares Go vs Swift IPA packaging
func BenchmarkIPAPackaging(b *testing.B) {
	if runtime.GOOS != "darwin" {
		b.Skip("Swift helpers only available on macOS")
	}

	_, swiftAvailable := findHelper(IPAPackBinary)

	// Create test app bundle with varying sizes
	tempDir := b.TempDir()
	sizes := []int{1, 10} // MB of app content

	for _, sizeMB := range sizes {
		appDir := filepath.Join(tempDir, fmt.Sprintf("TestApp-%dMB.app", sizeMB))
		payloadDir := filepath.Join(appDir, "Payload")
		_ = os.MkdirAll(payloadDir, 0o755)

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
		_ = os.WriteFile(filepath.Join(payloadDir, "Info.plist"), []byte(infoPlist), 0o644)

		// Add test content
		testContent := make([]byte, sizeMB*1024*1024)
		_ = os.WriteFile(filepath.Join(payloadDir, "content.bin"), testContent, 0o644)

		ipaPath := filepath.Join(tempDir, fmt.Sprintf("TestApp-%dMB.ipa", sizeMB))

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Benchmark Swift IPA packaging
		if swiftAvailable == nil {
			b.Run(fmt.Sprintf("Swift_%dMB", sizeMB), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					os.Remove(ipaPath)
					_, err := PackIPA(ctx, appDir, ipaPath, 6)
					if err != nil {
						b.Fatalf("Swift IPA packaging failed: %v", err)
					}
				}
			})
		}

		// Benchmark Go IPA packaging (basic zip)
		b.Run(fmt.Sprintf("Go_%dMB", sizeMB), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				os.Remove(ipaPath)
				err := createIPAGo(appDir, ipaPath)
				if err != nil {
					b.Fatalf("Go IPA packaging failed: %v", err)
				}
			}
		})
	}
}

// generateJWTGo generates a JWT using golang-jwt library
func generateJWTGo(keyID, issuerID string, privateKey *ecdsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": issuerID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 20).Unix(),
		"aud": "appstoreconnect-v1",
	})
	token.Header["kid"] = keyID
	token.Header["alg"] = "ES256"
	token.Header["typ"] = "JWT"

	return token.SignedString(privateKey)
}

// extractZipGo extracts a ZIP file using Go's standard library
func extractZipGo(zipPath, destDir string) error {
	cmd := exec.Command("unzip", "-o", zipPath, "-d", destDir)
	return cmd.Run()
}

// createIPAGo creates an IPA file using Go (wraps zip command)
func createIPAGo(appPath, ipaPath string) error {
	// IPA is essentially a ZIP with Payload directory
	cmd := exec.Command("zip", "-r", ipaPath, filepath.Base(appPath))
	cmd.Dir = filepath.Dir(appPath)
	return cmd.Run()
}

// createTestPNG creates a simple test PNG file
func createTestPNG(path string, width, height int) error {
	// Create a simple RGBA image and save as PNG
	// For benchmark purposes, we'll use a simple approach
	cmd := exec.Command("sips", "-s", "format", "png",
		"-Z", fmt.Sprintf("%d", max(width, height)),
		"/System/Library/CoreServices/DefaultDesktop.heic",
		"--out", path)
	return cmd.Run()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
