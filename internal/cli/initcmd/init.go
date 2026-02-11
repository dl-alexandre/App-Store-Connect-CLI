package initcmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/docs"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// InitCommand returns the root init command.
func InitCommand() *ffcli.Command {
	fs := flag.NewFlagSet("init", flag.ExitOnError)

	path := fs.String("path", "", "Output path for ASC.md (default: repo root or current directory)")
	force := fs.Bool("force", false, "Overwrite existing ASC.md")
	link := fs.Bool("link", true, "Update AGENTS.md and CLAUDE.md to reference ASC.md")

	return &ffcli.Command{
		Name:       "init",
		ShortUsage: "asc init [flags]",
		ShortHelp:  "Initialize ASC helper docs in the current repo.",
		LongHelp: `Initialize ASC helper docs in the current repo.

Examples:
  asc init
  asc init --path ./ASC.md
  asc init --force --link=false`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			result, err := docs.InitReference(docs.InitOptions{
				Path:  *path,
				Force: *force,
				Link:  *link,
			})
			if err != nil {
				return fmt.Errorf("init: %w", err)
			}
			return asc.PrintJSON(result)
		},
	}
}
