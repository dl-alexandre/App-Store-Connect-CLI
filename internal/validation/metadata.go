package validation

import (
	"fmt"
	"unicode/utf8"
)

func metadataLengthChecks(versionLocs []VersionLocalization, appInfoLocs []AppInfoLocalization) []CheckResult {
	var checks []CheckResult

	for _, loc := range versionLocs {
		if utf8.RuneCountInString(loc.Description) > LimitDescription {
			checks = append(checks, CheckResult{
				ID:           "metadata.length.description",
				Severity:     SeverityError,
				Locale:       loc.Locale,
				Field:        "description",
				ResourceType: "appStoreVersionLocalization",
				ResourceID:   loc.ID,
				Message:      fmt.Sprintf("description exceeds %d characters", LimitDescription),
				Remediation:  fmt.Sprintf("Shorten description to %d characters or fewer", LimitDescription),
			})
		}
		if utf8.RuneCountInString(loc.Keywords) > LimitKeywords {
			checks = append(checks, CheckResult{
				ID:           "metadata.length.keywords",
				Severity:     SeverityError,
				Locale:       loc.Locale,
				Field:        "keywords",
				ResourceType: "appStoreVersionLocalization",
				ResourceID:   loc.ID,
				Message:      fmt.Sprintf("keywords exceed %d characters", LimitKeywords),
				Remediation:  fmt.Sprintf("Shorten keywords to %d characters or fewer", LimitKeywords),
			})
		}
		if utf8.RuneCountInString(loc.WhatsNew) > LimitWhatsNew {
			checks = append(checks, CheckResult{
				ID:           "metadata.length.whats_new",
				Severity:     SeverityError,
				Locale:       loc.Locale,
				Field:        "whatsNew",
				ResourceType: "appStoreVersionLocalization",
				ResourceID:   loc.ID,
				Message:      fmt.Sprintf("what's new exceeds %d characters", LimitWhatsNew),
				Remediation:  fmt.Sprintf("Shorten what's new to %d characters or fewer", LimitWhatsNew),
			})
		}
		if utf8.RuneCountInString(loc.PromotionalText) > LimitPromotionalText {
			checks = append(checks, CheckResult{
				ID:           "metadata.length.promotional_text",
				Severity:     SeverityError,
				Locale:       loc.Locale,
				Field:        "promotionalText",
				ResourceType: "appStoreVersionLocalization",
				ResourceID:   loc.ID,
				Message:      fmt.Sprintf("promotional text exceeds %d characters", LimitPromotionalText),
				Remediation:  fmt.Sprintf("Shorten promotional text to %d characters or fewer", LimitPromotionalText),
			})
		}
	}

	for _, loc := range appInfoLocs {
		if utf8.RuneCountInString(loc.Name) > LimitName {
			checks = append(checks, CheckResult{
				ID:           "metadata.length.name",
				Severity:     SeverityError,
				Locale:       loc.Locale,
				Field:        "name",
				ResourceType: "appInfoLocalization",
				ResourceID:   loc.ID,
				Message:      fmt.Sprintf("name exceeds %d characters", LimitName),
				Remediation:  fmt.Sprintf("Shorten name to %d characters or fewer", LimitName),
			})
		}
		if utf8.RuneCountInString(loc.Subtitle) > LimitSubtitle {
			checks = append(checks, CheckResult{
				ID:           "metadata.length.subtitle",
				Severity:     SeverityError,
				Locale:       loc.Locale,
				Field:        "subtitle",
				ResourceType: "appInfoLocalization",
				ResourceID:   loc.ID,
				Message:      fmt.Sprintf("subtitle exceeds %d characters", LimitSubtitle),
				Remediation:  fmt.Sprintf("Shorten subtitle to %d characters or fewer", LimitSubtitle),
			})
		}
	}

	return checks
}
