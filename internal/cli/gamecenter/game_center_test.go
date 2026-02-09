package gamecenter

import "testing"

func TestGameCenterCommandConstructors(t *testing.T) {
	top := GameCenterCommand()
	if top == nil {
		t.Fatal("expected game-center command")
	}
	if top.Name == "" {
		t.Fatal("expected command name")
	}
	if len(top.Subcommands) == 0 {
		t.Fatal("expected subcommands")
	}

	if got := Command(); got == nil {
		t.Fatal("expected Command wrapper to return command")
	}

	constructors := []func() interface{}{
		func() interface{} { return GameCenterAchievementsCommand() },
		func() interface{} { return GameCenterLeaderboardsCommand() },
		func() interface{} { return GameCenterLeaderboardSetsCommand() },
		func() interface{} { return GameCenterGroupsCommand() },
		func() interface{} { return GameCenterDetailsCommand() },
		func() interface{} { return GameCenterAppVersionsCommand() },
		func() interface{} { return GameCenterEnabledVersionsCommand() },
		func() interface{} { return GameCenterMatchmakingCommand() },
		func() interface{} { return GameCenterChallengesCommand() },
		func() interface{} { return GameCenterActivitiesCommand() },
		func() interface{} { return GameCenterAchievementsV2Command() },
		func() interface{} { return GameCenterLeaderboardsV2Command() },
		func() interface{} { return GameCenterLeaderboardSetsV2Command() },
		func() interface{} { return GameCenterLeaderboardSetImagesCommand() },
	}
	for _, ctor := range constructors {
		if got := ctor(); got == nil {
			t.Fatal("expected constructor to return command")
		}
	}
}
