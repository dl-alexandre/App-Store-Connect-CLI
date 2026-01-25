package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

// CategoriesCommand returns the categories command with subcommands.
func CategoriesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("categories", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "categories",
		ShortUsage: "asc categories <subcommand> [flags]",
		ShortHelp:  "Manage App Store categories.",
		LongHelp: `Manage App Store categories.

Examples:
  asc categories list
  asc categories set --app APP_ID --primary GAMES`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			CategoriesListCommand(),
			CategoriesSetCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// CategoriesListCommand returns the categories list subcommand.
func CategoriesListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("categories list", flag.ExitOnError)

	limit := fs.Int("limit", 200, "Maximum results to fetch (1-200)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc categories list [flags]",
		ShortHelp:  "List available App Store categories.",
		LongHelp: `List available App Store categories.

Category IDs can be used when updating app information to set primary
and secondary categories.

Examples:
  asc categories list
  asc categories list --output table`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit < 1 || *limit > 200 {
				return fmt.Errorf("categories list: --limit must be between 1 and 200")
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("categories list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			categories, err := client.GetAppCategories(requestCtx, asc.WithAppCategoriesLimit(*limit))
			if err != nil {
				return fmt.Errorf("categories list: %w", err)
			}

			return printOutput(categories, *output, *pretty)
		},
	}
}

// CategoriesSetCommand returns the categories set subcommand.
func CategoriesSetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("categories set", flag.ExitOnError)

	appID := fs.String("app", os.Getenv("ASC_APP_ID"), "App ID (required)")
	primary := fs.String("primary", "", "Primary category ID (required)")
	secondary := fs.String("secondary", "", "Secondary category ID (optional)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "set",
		ShortUsage: "asc categories set --app APP_ID --primary CATEGORY_ID [--secondary CATEGORY_ID]",
		ShortHelp:  "Set primary and secondary categories for an app.",
		LongHelp: `Set the primary and secondary categories for an app.

Use 'asc categories list' to find valid category IDs.

Note: The app must have an editable version in PREPARE_FOR_SUBMISSION state.

Examples:
  asc categories set --app 123456789 --primary GAMES
  asc categories set --app 123456789 --primary GAMES --secondary ENTERTAINMENT
  asc categories set --app 123456789 --primary PHOTO_AND_VIDEO`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			appIDValue := strings.TrimSpace(*appID)
			primaryValue := strings.TrimSpace(*primary)
			secondaryValue := strings.TrimSpace(*secondary)

			if appIDValue == "" {
				return fmt.Errorf("categories set: --app is required")
			}
			if primaryValue == "" {
				return fmt.Errorf("categories set: --primary is required")
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("categories set: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get the current app info ID
			appInfos, err := client.GetAppInfos(requestCtx, appIDValue)
			if err != nil {
				return fmt.Errorf("categories set: failed to get app info: %w", err)
			}

			if len(appInfos.Data) == 0 {
				return fmt.Errorf("categories set: no app info found for app %s", appIDValue)
			}

			// Use the first (most recent) app info
			appInfoID := appInfos.Data[0].ID

			// Update categories
			resp, err := client.UpdateAppInfoCategories(requestCtx, appInfoID, primaryValue, secondaryValue)
			if err != nil {
				return fmt.Errorf("categories set: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}
