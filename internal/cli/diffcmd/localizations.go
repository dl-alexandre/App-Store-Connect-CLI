package diffcmd

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

type localizationDiffEndpoint struct {
	Kind      string `json:"kind"`
	Path      string `json:"path,omitempty"`
	VersionID string `json:"versionId,omitempty"`
}

type localizationDiffItem struct {
	Key    string `json:"key"`
	Locale string `json:"locale"`
	Field  string `json:"field"`
	Reason string `json:"reason"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
}

type localizationDiffPlan struct {
	Scope   string                   `json:"scope"`
	AppID   string                   `json:"appId"`
	Source  localizationDiffEndpoint `json:"source"`
	Target  localizationDiffEndpoint `json:"target"`
	Adds    []localizationDiffItem   `json:"adds"`
	Updates []localizationDiffItem   `json:"updates"`
	Deletes []localizationDiffItem   `json:"deletes"`
}

// DiffLocalizationsCommand compares localization metadata between two sources.
func DiffLocalizationsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("localizations", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (required, or ASC_APP_ID env)")
	path := fs.String("path", "", "Local .strings directory or file (source)")
	fromVersion := fs.String("from-version", "", "Remote source app store version ID")
	version := fs.String("version", "", "Remote target app store version ID (when using --path)")
	toVersion := fs.String("to-version", "", "Remote target app store version ID (when using --from-version)")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "localizations",
		ShortUsage: "asc diff localizations [flags]",
		ShortHelp:  "Diff localization metadata from local files or remote versions.",
		LongHelp: `Diff localization metadata from local files or remote versions.

Modes:
  Local vs remote:
    asc diff localizations --app "APP_ID" --path "./metadata/localizations" --version "VERSION_ID"

  Remote vs remote:
    asc diff localizations --app "APP_ID" --from-version "VERSION_ID_A" --to-version "VERSION_ID_B"`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return shared.UsageError("diff localizations does not accept positional arguments")
			}

			resolvedAppID := shared.ResolveAppID(*appID)
			if resolvedAppID == "" {
				return shared.UsageError("--app is required (or set ASC_APP_ID)")
			}

			sourcePath := strings.TrimSpace(*path)
			sourceVersion := strings.TrimSpace(*fromVersion)
			targetVersion := strings.TrimSpace(*version)
			targetToVersion := strings.TrimSpace(*toVersion)

			hasPath := sourcePath != ""
			hasFromVersion := sourceVersion != ""
			if hasPath && hasFromVersion {
				return shared.UsageError("--path and --from-version are mutually exclusive")
			}
			if !hasPath && !hasFromVersion {
				return shared.UsageError("either --path or --from-version is required")
			}

			var plan localizationDiffPlan
			var sourceValues map[string]map[string]string
			var targetValues map[string]map[string]string

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			if hasPath {
				if targetVersion == "" {
					return shared.UsageError("--version is required when using --path")
				}
				if targetToVersion != "" {
					return shared.UsageError("--to-version cannot be used with --path; use --version")
				}

				localValues, err := readAndValidateLocalLocalizations(sourcePath)
				if err != nil {
					return shared.UsageError(err.Error())
				}

				client, err := shared.GetASCClient()
				if err != nil {
					return fmt.Errorf("diff localizations: %w", err)
				}

				remoteValues, err := fetchVersionLocalizations(requestCtx, client, targetVersion)
				if err != nil {
					return fmt.Errorf("diff localizations: %w", err)
				}

				sourceValues = localValues
				targetValues = remoteValues
				plan = buildLocalizationDiffPlan(
					resolvedAppID,
					localizationDiffEndpoint{Kind: "local", Path: sourcePath},
					localizationDiffEndpoint{Kind: "remote", VersionID: targetVersion},
					sourceValues,
					targetValues,
				)
			} else {
				if targetToVersion == "" {
					return shared.UsageError("--to-version is required when using --from-version")
				}
				if targetVersion != "" {
					return shared.UsageError("--version cannot be used with --from-version; use --to-version")
				}

				client, err := shared.GetASCClient()
				if err != nil {
					return fmt.Errorf("diff localizations: %w", err)
				}

				fromValues, err := fetchVersionLocalizations(requestCtx, client, sourceVersion)
				if err != nil {
					return fmt.Errorf("diff localizations: %w", err)
				}
				toValues, err := fetchVersionLocalizations(requestCtx, client, targetToVersion)
				if err != nil {
					return fmt.Errorf("diff localizations: %w", err)
				}

				sourceValues = fromValues
				targetValues = toValues
				plan = buildLocalizationDiffPlan(
					resolvedAppID,
					localizationDiffEndpoint{Kind: "remote", VersionID: sourceVersion},
					localizationDiffEndpoint{Kind: "remote", VersionID: targetToVersion},
					sourceValues,
					targetValues,
				)
			}

			return shared.PrintOutputWithRenderers(
				plan,
				*output.Output,
				*output.Pretty,
				func() error {
					renderLocalizationDiffTable(plan)
					return nil
				},
				func() error {
					renderLocalizationDiffMarkdown(plan)
					return nil
				},
			)
		},
	}
}

func readAndValidateLocalLocalizations(inputPath string) (map[string]map[string]string, error) {
	valuesByLocale, err := shared.ReadLocalizationStrings(inputPath, nil)
	if err != nil {
		return nil, err
	}

	normalized := make(map[string]map[string]string, len(valuesByLocale))
	for locale, values := range valuesByLocale {
		if err := validateVersionLocalizationFields(locale, values); err != nil {
			return nil, err
		}
		normalized[locale] = normalizeLocalizationValues(values)
	}

	return normalized, nil
}

func validateVersionLocalizationFields(locale string, values map[string]string) error {
	fields := shared.VersionLocalizationKeys()
	allowed := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		allowed[field] = struct{}{}
	}

	unknown := make([]string, 0)
	for key := range values {
		if _, ok := allowed[key]; !ok {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) == 0 {
		return nil
	}

	sort.Strings(unknown)
	return fmt.Errorf("unsupported keys for locale %q: %s", locale, strings.Join(unknown, ", "))
}

func fetchVersionLocalizations(ctx context.Context, client *asc.Client, versionID string) (map[string]map[string]string, error) {
	firstPage, err := client.GetAppStoreVersionLocalizations(
		ctx,
		versionID,
		asc.WithAppStoreVersionLocalizationsLimit(200),
	)
	if err != nil {
		return nil, err
	}

	resp := firstPage
	if firstPage != nil && firstPage.Links.Next != "" {
		paginated, err := asc.PaginateAll(
			ctx,
			firstPage,
			func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
				return client.GetAppStoreVersionLocalizations(
					ctx,
					"",
					asc.WithAppStoreVersionLocalizationsNextURL(nextURL),
				)
			},
		)
		if err != nil {
			return nil, err
		}

		typed, ok := paginated.(*asc.AppStoreVersionLocalizationsResponse)
		if !ok {
			return nil, fmt.Errorf("unexpected pagination response type")
		}
		resp = typed
	}

	valuesByLocale := make(map[string]map[string]string)
	if resp == nil {
		return valuesByLocale, nil
	}

	for _, item := range resp.Data {
		locale := strings.TrimSpace(item.Attributes.Locale)
		if locale == "" {
			continue
		}
		if _, exists := valuesByLocale[locale]; exists {
			return nil, fmt.Errorf("duplicate locale %q in remote version %q", locale, versionID)
		}
		valuesByLocale[locale] = normalizeLocalizationValues(shared.MapVersionLocalizationStrings(item.Attributes))
	}

	return valuesByLocale, nil
}

func normalizeLocalizationValues(values map[string]string) map[string]string {
	normalized := make(map[string]string, len(values))
	for key, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		normalized[key] = trimmed
	}
	return normalized
}

func buildLocalizationDiffPlan(
	appID string,
	source localizationDiffEndpoint,
	target localizationDiffEndpoint,
	sourceValues map[string]map[string]string,
	targetValues map[string]map[string]string,
) localizationDiffPlan {
	plan := localizationDiffPlan{
		Scope:   "localizations",
		AppID:   appID,
		Source:  source,
		Target:  target,
		Adds:    make([]localizationDiffItem, 0),
		Updates: make([]localizationDiffItem, 0),
		Deletes: make([]localizationDiffItem, 0),
	}

	localesMap := make(map[string]struct{})
	for locale := range sourceValues {
		localesMap[locale] = struct{}{}
	}
	for locale := range targetValues {
		localesMap[locale] = struct{}{}
	}

	locales := make([]string, 0, len(localesMap))
	for locale := range localesMap {
		locales = append(locales, locale)
	}
	sort.Strings(locales)
	fields := shared.VersionLocalizationKeys()

	for _, locale := range locales {
		sourceFields := sourceValues[locale]
		targetFields := targetValues[locale]

		for _, field := range fields {
			sourceValue, sourceOK := sourceFields[field]
			targetValue, targetOK := targetFields[field]
			key := fmt.Sprintf("%s:%s", locale, field)

			switch {
			case !sourceOK && targetOK:
				plan.Adds = append(plan.Adds, localizationDiffItem{
					Key:    key,
					Locale: locale,
					Field:  field,
					Reason: "field exists in target but not in source",
					To:     targetValue,
				})
			case sourceOK && !targetOK:
				plan.Deletes = append(plan.Deletes, localizationDiffItem{
					Key:    key,
					Locale: locale,
					Field:  field,
					Reason: "field exists in source but not in target",
					From:   sourceValue,
				})
			case sourceOK && targetOK && sourceValue != targetValue:
				plan.Updates = append(plan.Updates, localizationDiffItem{
					Key:    key,
					Locale: locale,
					Field:  field,
					Reason: "field value differs",
					From:   sourceValue,
					To:     targetValue,
				})
			}
		}
	}

	return plan
}

func renderLocalizationDiffTable(plan localizationDiffPlan) {
	headers := []string{"change", "key", "locale", "field", "reason", "from", "to"}
	asc.RenderTable(headers, buildLocalizationDiffRows(plan))
}

func renderLocalizationDiffMarkdown(plan localizationDiffPlan) {
	headers := []string{"change", "key", "locale", "field", "reason", "from", "to"}
	asc.RenderMarkdown(headers, buildLocalizationDiffRows(plan))
}

func buildLocalizationDiffRows(plan localizationDiffPlan) [][]string {
	rows := make([][]string, 0, len(plan.Adds)+len(plan.Updates)+len(plan.Deletes))

	for _, item := range plan.Adds {
		rows = append(rows, []string{
			"add",
			item.Key,
			item.Locale,
			item.Field,
			item.Reason,
			"",
			sanitizeDiffCell(item.To),
		})
	}
	for _, item := range plan.Updates {
		rows = append(rows, []string{
			"update",
			item.Key,
			item.Locale,
			item.Field,
			item.Reason,
			sanitizeDiffCell(item.From),
			sanitizeDiffCell(item.To),
		})
	}
	for _, item := range plan.Deletes {
		rows = append(rows, []string{
			"delete",
			item.Key,
			item.Locale,
			item.Field,
			item.Reason,
			sanitizeDiffCell(item.From),
			"",
		})
	}

	if len(rows) == 0 {
		rows = append(rows, []string{"none", "", "", "", "no changes", "", ""})
	}
	return rows
}

func sanitizeDiffCell(value string) string {
	normalized := strings.ReplaceAll(value, "\n", "\\n")
	const maxLen = 80
	if len(normalized) <= maxLen {
		return normalized
	}
	return normalized[:maxLen-3] + "..."
}
