package update

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallScriptUsesReleaseAssetNaming(t *testing.T) {
	path := filepath.Join("..", "..", "install.sh")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read install.sh: %v", err)
	}
	script := string(data)

	checks := []string{
		"releases/latest",
		"releases/download/${VERSION}",
		"Darwin) OS=\"macOS\"",
		"Linux) OS=\"linux\"",
		"ASSET=\"${BIN_NAME}_${VERSION}_${OS}_${ARCH}\"",
		"CHECKSUMS_ASSET=\"${BIN_NAME}_${VERSION}_checksums.txt\"",
	}

	for _, snippet := range checks {
		if !strings.Contains(script, snippet) {
			t.Errorf("install.sh missing %q", snippet)
		}
	}
}
