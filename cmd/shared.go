package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// Deprecated: use CleanupTempPrivateKeys to remove all tracked temp keys.
func CleanupTempPrivateKey() {
	shared.CleanupTempPrivateKey()
}

func CleanupTempPrivateKeys() {
	shared.CleanupTempPrivateKeys()
}

func Bold(s string) string {
	return shared.Bold(s)
}

func DefaultUsageFunc(c *ffcli.Command) string {
	return shared.DefaultUsageFunc(c)
}
