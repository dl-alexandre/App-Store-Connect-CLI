package asc

import (
	"fmt"
	"strings"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/validation"
)

func init() {
	registerDirect(func(v *validation.Report, render func([]string, [][]string)) error {
		h, r := validationSummaryRows(v)
		render(h, r)
		oh, or := validationCheckRows(v)
		render(oh, or)
		return nil
	})
}

func validationSummaryRows(report *validation.Report) ([]string, [][]string) {
	headers := []string{"App ID", "Version ID", "Version", "Platform", "Errors", "Warnings", "Infos", "Blocking", "Strict"}
	rows := [][]string{{
		report.AppID,
		report.VersionID,
		report.VersionString,
		report.Platform,
		fmt.Sprintf("%d", report.Summary.Errors),
		fmt.Sprintf("%d", report.Summary.Warnings),
		fmt.Sprintf("%d", report.Summary.Infos),
		fmt.Sprintf("%d", report.Summary.Blocking),
		formatBool(report.Strict),
	}}
	return headers, rows
}

func validationCheckRows(report *validation.Report) ([]string, [][]string) {
	headers := []string{"Severity", "Check ID", "Locale", "Field", "Resource", "Message", "Remediation"}
	if report == nil || len(report.Checks) == 0 {
		return headers, [][]string{{"info", "validation.ok", "", "", "", "No issues found", ""}}
	}

	rows := make([][]string, 0, len(report.Checks))
	for _, check := range report.Checks {
		rows = append(rows, []string{
			string(check.Severity),
			check.ID,
			check.Locale,
			check.Field,
			formatResource(check.ResourceType, check.ResourceID),
			check.Message,
			check.Remediation,
		})
	}
	return headers, rows
}

func formatResource(resourceType, resourceID string) string {
	if resourceType == "" && resourceID == "" {
		return ""
	}
	if resourceID == "" {
		return resourceType
	}
	if resourceType == "" {
		return resourceID
	}
	return strings.TrimSpace(resourceType) + ":" + strings.TrimSpace(resourceID)
}

func formatBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
