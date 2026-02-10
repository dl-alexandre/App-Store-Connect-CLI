package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"testing"
)

func TestGameCenterMemberLocalizationsCreateValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing leaderboard-set-id",
			args: []string{"game-center", "leaderboard-sets", "member-localizations", "create", "--leaderboard-id", "LB_ID", "--locale", "en-US", "--name", "Test"},
		},
		{
			name: "missing leaderboard-id",
			args: []string{"game-center", "leaderboard-sets", "member-localizations", "create", "--leaderboard-set-id", "SET_ID", "--locale", "en-US", "--name", "Test"},
		},
		{
			name: "missing locale",
			args: []string{"game-center", "leaderboard-sets", "member-localizations", "create", "--leaderboard-set-id", "SET_ID", "--leaderboard-id", "LB_ID", "--name", "Test"},
		},
		{
			name: "missing name",
			args: []string{"game-center", "leaderboard-sets", "member-localizations", "create", "--leaderboard-set-id", "SET_ID", "--leaderboard-id", "LB_ID", "--locale", "en-US"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			stdout, _ := captureOutput(t, func() {
				if err := root.Parse(test.args); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				err := root.Run(context.Background())
				if !errors.Is(err, flag.ErrHelp) {
					t.Fatalf("expected ErrHelp, got %v", err)
				}
			})

			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
		})
	}
}

func TestGameCenterMemberLocalizationsUpdateValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing id",
			args: []string{"game-center", "leaderboard-sets", "member-localizations", "update", "--name", "New Name"},
		},
		{
			name: "no update flags",
			args: []string{"game-center", "leaderboard-sets", "member-localizations", "update", "--id", "LOCALIZATION_ID"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			stdout, _ := captureOutput(t, func() {
				if err := root.Parse(test.args); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				err := root.Run(context.Background())
				if !errors.Is(err, flag.ErrHelp) {
					t.Fatalf("expected ErrHelp, got %v", err)
				}
			})

			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
		})
	}
}

func TestGameCenterMemberLocalizationsDeleteValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing id",
			args: []string{"game-center", "leaderboard-sets", "member-localizations", "delete", "--confirm"},
		},
		{
			name: "missing confirm",
			args: []string{"game-center", "leaderboard-sets", "member-localizations", "delete", "--id", "LOCALIZATION_ID"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			stdout, _ := captureOutput(t, func() {
				if err := root.Parse(test.args); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				err := root.Run(context.Background())
				if !errors.Is(err, flag.ErrHelp) {
					t.Fatalf("expected ErrHelp, got %v", err)
				}
			})

			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
		})
	}
}
