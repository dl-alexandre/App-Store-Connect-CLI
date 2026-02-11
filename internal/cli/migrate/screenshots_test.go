package migrate

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestInferScreenshotDisplayType_FromFilenameAndDimensions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "iphone_65_screen.png")
	writePNG(t, path, 1242, 2688)

	displayType, err := inferScreenshotDisplayType(path)
	if err != nil {
		t.Fatalf("inferScreenshotDisplayType() error: %v", err)
	}
	if displayType != "APP_IPHONE_65" {
		t.Fatalf("expected APP_IPHONE_65, got %q", displayType)
	}
}

func TestInferScreenshotDisplayType_FromDimensionsOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "screen.png")
	writePNG(t, path, 1242, 2688)

	displayType, err := inferScreenshotDisplayType(path)
	if err != nil {
		t.Fatalf("inferScreenshotDisplayType() error: %v", err)
	}
	if displayType != "APP_IPHONE_65" {
		t.Fatalf("expected APP_IPHONE_65, got %q", displayType)
	}
}

func TestInferScreenshotDisplayType_MapsIPhone69AliasToAppIPhone67(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "iphone_69_screen.png")
	writePNG(t, path, 1290, 2796)

	displayType, err := inferScreenshotDisplayType(path)
	if err != nil {
		t.Fatalf("inferScreenshotDisplayType() error: %v", err)
	}
	if displayType != "APP_IPHONE_67" {
		t.Fatalf("expected APP_IPHONE_67, got %q", displayType)
	}
}

func TestInferScreenshotDisplayType_AcceptsLatestLargeIPhoneDimensions(t *testing.T) {
	testCases := []struct {
		name   string
		width  int
		height int
	}{
		{name: "1320x2868", width: 1320, height: 2868},
		{name: "1260x2736", width: 1260, height: 2736},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "screen.png")
			writePNG(t, path, tc.width, tc.height)

			displayType, err := inferScreenshotDisplayType(path)
			if err != nil {
				t.Fatalf("inferScreenshotDisplayType() error: %v", err)
			}
			if displayType != "APP_IPHONE_67" {
				t.Fatalf("expected APP_IPHONE_67, got %q", displayType)
			}
		})
	}
}

func TestInferScreenshotDisplayType_UnknownSize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "screen.png")
	writePNG(t, path, 120, 240)

	_, err := inferScreenshotDisplayType(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDiscoverScreenshotPlan_NormalizesLocale(t *testing.T) {
	root := t.TempDir()
	localeDir := filepath.Join(root, "en_US")
	if err := os.MkdirAll(localeDir, 0o755); err != nil {
		t.Fatalf("mkdir locale dir: %v", err)
	}
	writePNG(t, filepath.Join(localeDir, "iphone_65_screen.png"), 1242, 2688)

	plans, _, err := discoverScreenshotPlan(root)
	if err != nil {
		t.Fatalf("discoverScreenshotPlan() error: %v", err)
	}
	if len(plans) != 1 {
		t.Fatalf("expected 1 plan, got %d", len(plans))
	}
	if plans[0].Locale != "en-US" {
		t.Fatalf("expected locale en-US, got %q", plans[0].Locale)
	}
	if plans[0].DisplayType != "APP_IPHONE_65" {
		t.Fatalf("expected display type APP_IPHONE_65, got %q", plans[0].DisplayType)
	}
	if len(plans[0].Files) != 1 {
		t.Fatalf("expected 1 screenshot file, got %d", len(plans[0].Files))
	}
}

func writePNG(t *testing.T, path string, width, height int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 10, G: 20, B: 30, A: 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create png: %v", err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
}
