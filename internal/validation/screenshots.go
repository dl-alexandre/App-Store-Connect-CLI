package validation

import (
	"fmt"
	"strings"
)

type screenshotSize struct {
	Width  int
	Height int
}

var screenshotSizeCatalog = map[string][]screenshotSize{
	"APP_IPHONE_69": {
		{Width: 1260, Height: 2736},
		{Width: 1290, Height: 2796},
		{Width: 1320, Height: 2868},
		{Width: 1284, Height: 2778},
	},
	"APP_IPHONE_67": {
		{Width: 1260, Height: 2736},
		{Width: 1290, Height: 2796},
		{Width: 1320, Height: 2868},
		{Width: 1284, Height: 2778},
	},
	"APP_IPHONE_65":         {{Width: 1242, Height: 2688}},
	"APP_IPHONE_61":         {{Width: 1179, Height: 2556}, {Width: 1170, Height: 2532}},
	"APP_IPHONE_58":         {{Width: 1125, Height: 2436}},
	"APP_IPHONE_55":         {{Width: 1242, Height: 2208}},
	"APP_IPHONE_47":         {{Width: 750, Height: 1334}},
	"APP_IPHONE_40":         {{Width: 640, Height: 1136}},
	"APP_IPHONE_35":         {{Width: 640, Height: 960}},
	"APP_IPAD_PRO_3GEN_129": {{Width: 2048, Height: 2732}},
	"APP_IPAD_PRO_3GEN_11":  {{Width: 1668, Height: 2388}},
	"APP_IPAD_PRO_129":      {{Width: 2048, Height: 2732}},
	"APP_IPAD_105":          {{Width: 1668, Height: 2224}},
	"APP_IPAD_97":           {{Width: 1536, Height: 2048}},
	"APP_DESKTOP": {
		{Width: 1280, Height: 800},
		{Width: 1440, Height: 900},
		{Width: 2560, Height: 1600},
		{Width: 2880, Height: 1800},
	},
	"APP_APPLE_TV":         {{Width: 1920, Height: 1080}, {Width: 3840, Height: 2160}},
	"APP_APPLE_VISION_PRO": {{Width: 3840, Height: 2160}},
	"APP_WATCH_ULTRA":      {{Width: 422, Height: 514}, {Width: 410, Height: 502}},
	"APP_WATCH_SERIES_10":  {{Width: 416, Height: 496}},
	"APP_WATCH_SERIES_7":   {{Width: 396, Height: 484}},
	"APP_WATCH_SERIES_4":   {{Width: 368, Height: 448}},
	"APP_WATCH_SERIES_3":   {{Width: 312, Height: 390}},
}

func screenshotChecks(platform string, sets []ScreenshotSet) []CheckResult {
	var checks []CheckResult
	normalizedPlatform := strings.ToUpper(strings.TrimSpace(platform))

	for _, set := range sets {
		displayType := strings.TrimSpace(set.DisplayType)
		if displayType == "" {
			continue
		}

		expectedPlatform := platformForDisplayType(displayType)
		if normalizedPlatform != "" && expectedPlatform != "" && normalizedPlatform != expectedPlatform {
			checks = append(checks, CheckResult{
				ID:           "screenshots.display_type_platform_mismatch",
				Severity:     SeverityError,
				Locale:       set.Locale,
				ResourceType: "appScreenshotSet",
				ResourceID:   set.ID,
				Message:      fmt.Sprintf("display type %s is not valid for platform %s", displayType, normalizedPlatform),
				Remediation:  "Use a screenshot display type compatible with the target platform",
			})
		}

		sizes := screenshotSizesForDisplayType(displayType)
		if len(sizes) == 0 {
			checks = append(checks, CheckResult{
				ID:           "screenshots.display_type_unknown",
				Severity:     SeverityWarning,
				Locale:       set.Locale,
				ResourceType: "appScreenshotSet",
				ResourceID:   set.ID,
				Message:      fmt.Sprintf("unknown screenshot display type %s", displayType),
				Remediation:  "Verify the display type and update the size catalog if needed",
			})
			continue
		}

		for _, shot := range set.Screenshots {
			if shot.Width <= 0 || shot.Height <= 0 {
				checks = append(checks, CheckResult{
					ID:           "screenshots.missing_dimensions",
					Severity:     SeverityWarning,
					Locale:       set.Locale,
					ResourceType: "appScreenshot",
					ResourceID:   shot.ID,
					Message:      fmt.Sprintf("missing screenshot dimensions for %s", shot.FileName),
					Remediation:  "Re-upload the screenshot so dimensions are available",
				})
				continue
			}

			if !matchesScreenshotSize(shot.Width, shot.Height, sizes) {
				checks = append(checks, CheckResult{
					ID:           "screenshots.dimension_mismatch",
					Severity:     SeverityError,
					Locale:       set.Locale,
					ResourceType: "appScreenshot",
					ResourceID:   shot.ID,
					Message:      fmt.Sprintf("screenshot size %dx%d does not match %s requirements", shot.Width, shot.Height, displayType),
					Remediation:  "Upload a screenshot with an approved size for the display type",
				})
			}
		}
	}

	return checks
}

func screenshotSizesForDisplayType(displayType string) []screenshotSize {
	if sizes, ok := screenshotSizeCatalog[displayType]; ok {
		return sizes
	}
	if strings.HasPrefix(displayType, "IMESSAGE_APP_") {
		base := strings.TrimPrefix(displayType, "IMESSAGE_APP_")
		if sizes, ok := screenshotSizeCatalog["APP_"+base]; ok {
			return sizes
		}
	}
	return nil
}

func matchesScreenshotSize(width, height int, sizes []screenshotSize) bool {
	for _, size := range sizes {
		if width == size.Width && height == size.Height {
			return true
		}
		if width == size.Height && height == size.Width {
			return true
		}
	}
	return false
}

func platformForDisplayType(displayType string) string {
	switch {
	case strings.HasPrefix(displayType, "APP_IPHONE"),
		strings.HasPrefix(displayType, "APP_IPAD"),
		strings.HasPrefix(displayType, "IMESSAGE_APP_"),
		strings.HasPrefix(displayType, "APP_WATCH"):
		return "IOS"
	case strings.HasPrefix(displayType, "APP_DESKTOP"):
		return "MAC_OS"
	case strings.HasPrefix(displayType, "APP_APPLE_TV"):
		return "TV_OS"
	case strings.HasPrefix(displayType, "APP_APPLE_VISION_PRO"):
		return "VISION_OS"
	default:
		return ""
	}
}
