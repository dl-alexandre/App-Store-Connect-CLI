package validation

import "testing"

func TestScreenshotChecks_Mismatch(t *testing.T) {
	sets := []ScreenshotSet{
		{
			ID:          "set-1",
			DisplayType: "APP_IPHONE_65",
			Locale:      "en-US",
			Screenshots: []Screenshot{
				{ID: "shot-1", FileName: "shot.png", Width: 100, Height: 100},
			},
		},
	}

	checks := screenshotChecks("IOS", sets)
	if !hasCheckID(checks, "screenshots.dimension_mismatch") {
		t.Fatalf("expected dimension mismatch check")
	}
}

func TestScreenshotChecks_Pass(t *testing.T) {
	sets := []ScreenshotSet{
		{
			ID:          "set-1",
			DisplayType: "APP_IPHONE_65",
			Locale:      "en-US",
			Screenshots: []Screenshot{
				{ID: "shot-1", FileName: "shot.png", Width: 1242, Height: 2688},
			},
		},
	}

	checks := screenshotChecks("IOS", sets)
	if len(checks) != 0 {
		t.Fatalf("expected no checks, got %d", len(checks))
	}
}

func TestScreenshotChecks_PassLatestLargeIPhoneSizes(t *testing.T) {
	sets := []ScreenshotSet{
		{
			ID:          "set-1",
			DisplayType: "APP_IPHONE_67",
			Locale:      "en-US",
			Screenshots: []Screenshot{
				{ID: "shot-1", FileName: "shot-1.png", Width: 1260, Height: 2736},
				{ID: "shot-2", FileName: "shot-2.png", Width: 1320, Height: 2868},
			},
		},
	}

	checks := screenshotChecks("IOS", sets)
	if len(checks) != 0 {
		t.Fatalf("expected no checks, got %d (%v)", len(checks), checks)
	}
}

func TestScreenshotChecks_PassDesktopAndWatchUltraNewestSizes(t *testing.T) {
	sets := []ScreenshotSet{
		{
			ID:          "set-mac",
			DisplayType: "APP_DESKTOP",
			Locale:      "en-US",
			Screenshots: []Screenshot{
				{ID: "shot-mac", FileName: "mac.png", Width: 2880, Height: 1800},
			},
		},
		{
			ID:          "set-watch",
			DisplayType: "APP_WATCH_ULTRA",
			Locale:      "en-US",
			Screenshots: []Screenshot{
				{ID: "shot-watch", FileName: "watch.png", Width: 422, Height: 514},
			},
		},
	}

	checks := screenshotChecks("IOS", sets)
	if len(checks) != 1 {
		t.Fatalf("expected one platform mismatch check for APP_DESKTOP under IOS, got %d (%v)", len(checks), checks)
	}
	if checks[0].ID != "screenshots.display_type_platform_mismatch" {
		t.Fatalf("expected platform mismatch check, got %s", checks[0].ID)
	}

	iosOnly := screenshotChecks("IOS", []ScreenshotSet{sets[1]})
	if len(iosOnly) != 0 {
		t.Fatalf("expected no checks for watch ultra IOS set, got %d (%v)", len(iosOnly), iosOnly)
	}

	macOnly := screenshotChecks("MAC_OS", []ScreenshotSet{sets[0]})
	if len(macOnly) != 0 {
		t.Fatalf("expected no checks for desktop MAC_OS set, got %d (%v)", len(macOnly), macOnly)
	}
}
