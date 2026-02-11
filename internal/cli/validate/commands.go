package validate

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// ValidateCommand returns the asc validate command.
func ValidateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID)")
	versionID := fs.String("version-id", "", "App Store version ID (required)")
	platform := fs.String("platform", "", "Platform: IOS, MAC_OS, TV_OS, VISION_OS")
	strict := fs.Bool("strict", false, "Treat warnings as errors (exit non-zero)")
	output := fs.String("output", shared.DefaultOutputFormat(), "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "validate",
		ShortUsage: "asc validate --app \"APP_ID\" --version-id \"VERSION_ID\" [flags]",
		ShortHelp:  "Validate metadata, screenshots, and age ratings before submission.",
		LongHelp: `Validate pre-submission readiness for an App Store version.

Checks:
  - Metadata length limits
  - Required fields and localizations
  - Screenshot size compatibility
  - Age rating completeness

Examples:
  asc validate --app "APP_ID" --version-id "VERSION_ID"
  asc validate --app "APP_ID" --version-id "VERSION_ID" --platform IOS --output table
  asc validate --app "APP_ID" --version-id "VERSION_ID" --strict`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if strings.TrimSpace(*versionID) == "" {
				fmt.Fprintln(os.Stderr, "Error: --version-id is required")
				return flag.ErrHelp
			}

			resolvedAppID := shared.ResolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			var normalizedPlatform string
			if strings.TrimSpace(*platform) != "" {
				value, err := shared.NormalizeAppStoreVersionPlatform(*platform)
				if err != nil {
					return fmt.Errorf("validate: %w", err)
				}
				normalizedPlatform = value
			}

			return runValidate(ctx, validateOptions{
				AppID:     resolvedAppID,
				VersionID: strings.TrimSpace(*versionID),
				Platform:  normalizedPlatform,
				Strict:    *strict,
				Output:    *output,
				Pretty:    *pretty,
			})
		},
	}
}
