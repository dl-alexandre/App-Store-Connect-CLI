package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

const promoCodesMaxLimit = 200
const promoCodesMaxQuantity = 10

// PromoCodesCommand returns the promo codes command with subcommands.
func PromoCodesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("promocodes", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "promocodes",
		ShortUsage: "asc promocodes <subcommand> [flags]",
		ShortHelp:  "Manage App Store Connect promo codes.",
		LongHelp: `Manage App Store Connect promo codes.

Promo codes can be generated for apps or subscriptions to distribute free access.

Examples:
  asc promocodes list --app "123456789"
  asc promocodes generate --app "123456789" --type app --quantity 5 --output "./promo-codes.txt"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			PromoCodesListCommand(),
			PromoCodesGenerateCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// PromoCodesListCommand returns the promo codes list subcommand.
func PromoCodesListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc promocodes list [flags]",
		ShortHelp:  "List promo codes for an app in App Store Connect.",
		LongHelp: `List promo codes for an app in App Store Connect.

Examples:
  asc promocodes list --app "123456789"
  asc promocodes list --app "123456789" --limit 10
  asc promocodes list --app "123456789" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > promoCodesMaxLimit) {
				return fmt.Errorf("promocodes list: --limit must be between 1 and %d", promoCodesMaxLimit)
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("promocodes list: %w", err)
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintf(os.Stderr, "Error: --app is required (or set ASC_APP_ID)\n\n")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("promocodes list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.PromoCodesOption{
				asc.WithPromoCodesLimit(*limit),
				asc.WithPromoCodesNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithPromoCodesLimit(promoCodesMaxLimit))
				firstPage, err := client.GetPromoCodes(requestCtx, resolvedAppID, paginateOpts...)
				if err != nil {
					return fmt.Errorf("promocodes list: failed to fetch: %w", err)
				}

				pages, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetPromoCodes(ctx, resolvedAppID, asc.WithPromoCodesNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("promocodes list: %w", err)
				}

				return printOutput(pages, *output, *pretty)
			}

			resp, err := client.GetPromoCodes(requestCtx, resolvedAppID, opts...)
			if err != nil {
				return fmt.Errorf("promocodes list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// PromoCodesGenerateCommand returns the promo codes generate subcommand.
func PromoCodesGenerateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (required, or ASC_APP_ID env)")
	codeType := fs.String("type", "", "Promo code type: app or subscription (required)")
	quantity := fs.Int("quantity", 0, "Number of promo codes to generate (1-10)")
	outputPath := fs.String("output", "", "Output file path for promo codes (one per line)")
	outputFormat := fs.String("output-format", "json", "Output format for metadata: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "generate",
		ShortUsage: "asc promocodes generate [flags]",
		ShortHelp:  "Generate new promo codes for an app.",
		LongHelp: `Generate new promo codes for an app.

Examples:
  asc promocodes generate --app "123456789" --type app --quantity 5
  asc promocodes generate --app "123456789" --type subscription --quantity 3 --output "./promo-codes.txt"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintf(os.Stderr, "Error: --app is required (or set ASC_APP_ID)\n\n")
				return flag.ErrHelp
			}

			codeTypeValue := strings.ToLower(strings.TrimSpace(*codeType))
			if codeTypeValue == "" {
				fmt.Fprintln(os.Stderr, "Error: --type is required")
				return flag.ErrHelp
			}

			if *quantity == 0 {
				fmt.Fprintln(os.Stderr, "Error: --quantity is required")
				return flag.ErrHelp
			}
			if *quantity < 1 || *quantity > promoCodesMaxQuantity {
				return fmt.Errorf("promocodes generate: --quantity must be between 1 and %d", promoCodesMaxQuantity)
			}

			var productType asc.PromoCodeProductType
			switch codeTypeValue {
			case "app":
				productType = asc.PromoCodeProductTypeApp
			case "subscription":
				productType = asc.PromoCodeProductTypeSubscription
			default:
				return fmt.Errorf("promocodes generate: --type must be app or subscription")
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("promocodes generate: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			req := asc.PromoCodeCreateRequest{
				Data: asc.PromoCodeCreateData{
					Type: asc.ResourceTypePromoCodes,
					Attributes: asc.PromoCodeCreateAttributes{
						ProductType: productType,
						Quantity:    *quantity,
					},
					Relationships: asc.PromoCodeCreateRelationships{
						App: asc.Relationship{
							Data: asc.ResourceData{
								Type: asc.ResourceTypeApps,
								ID:   resolvedAppID,
							},
						},
					},
				},
			}

			resp, err := client.CreatePromoCodes(requestCtx, req)
			if err != nil {
				return fmt.Errorf("promocodes generate: failed to generate: %w", err)
			}

			if strings.TrimSpace(*outputPath) != "" {
				codes := extractPromoCodes(resp)
				if len(codes) == 0 {
					return fmt.Errorf("promocodes generate: no codes returned to write")
				}
				if err := writePromoCodesFile(*outputPath, codes); err != nil {
					return fmt.Errorf("promocodes generate: %w", err)
				}
			}

			return printOutput(resp, *outputFormat, *pretty)
		},
	}
}

func extractPromoCodes(resp *asc.PromoCodesResponse) []string {
	if resp == nil {
		return nil
	}
	codes := make([]string, 0, len(resp.Data))
	for _, item := range resp.Data {
		code := strings.TrimSpace(item.Attributes.Code)
		if code == "" {
			continue
		}
		codes = append(codes, code)
	}
	return codes
}

func writePromoCodesFile(path string, codes []string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := openNewFileNoFollow(path, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("output file already exists: %w", err)
		}
		return err
	}
	defer file.Close()

	for _, code := range codes {
		trimmed := strings.TrimSpace(code)
		if trimmed == "" {
			continue
		}
		if _, err := fmt.Fprintln(file, trimmed); err != nil {
			return err
		}
	}
	return file.Sync()
}
