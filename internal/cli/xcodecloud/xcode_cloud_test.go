package xcodecloud

import "testing"

func TestXcodeCloudCommandConstructors(t *testing.T) {
	top := XcodeCloudCommand()
	if top == nil {
		t.Fatal("expected xcode-cloud command")
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
		func() interface{} { return XcodeCloudRunCommand() },
		func() interface{} { return XcodeCloudStatusCommand() },
		func() interface{} { return XcodeCloudWorkflowsCommand() },
		func() interface{} { return XcodeCloudBuildRunsCommand() },
		func() interface{} { return XcodeCloudActionsCommand() },
		func() interface{} { return XcodeCloudArtifactsCommand() },
		func() interface{} { return XcodeCloudTestResultsCommand() },
		func() interface{} { return XcodeCloudIssuesCommand() },
		func() interface{} { return XcodeCloudScmCommand() },
		func() interface{} { return XcodeCloudProductsCommand() },
		func() interface{} { return XcodeCloudMacOSVersionsCommand() },
		func() interface{} { return XcodeCloudXcodeVersionsCommand() },
	}
	for _, ctor := range constructors {
		if got := ctor(); got == nil {
			t.Fatal("expected constructor to return command")
		}
	}
}
