package syscli

import (
	"context"
	"runtime"
	"testing"
)

func TestIsDarwin(t *testing.T) {
	got := IsDarwin()
	want := runtime.GOOS == "darwin"
	if got != want {
		t.Errorf("IsDarwin() = %v, want %v", got, want)
	}
}

func TestListSimulators_NonDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("test only runs on non-darwin platforms")
	}
	_, err := ListSimulators(context.Background(), false)
	if err == nil {
		t.Error("expected error on non-darwin platform")
	}
}

func TestSimulatorScreenshot_NonDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("test only runs on non-darwin platforms")
	}
	err := SimulatorScreenshot(context.Background(), "fake-udid", "/tmp/out.png")
	if err == nil {
		t.Error("expected error on non-darwin platform")
	}
}

func TestSimulatorInstall_NonDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("test only runs on non-darwin platforms")
	}
	err := SimulatorInstall(context.Background(), "fake-udid", "/tmp/app.app")
	if err == nil {
		t.Error("expected error on non-darwin platform")
	}
}

func TestSimulatorLaunch_NonDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("test only runs on non-darwin platforms")
	}
	err := SimulatorLaunch(context.Background(), "fake-udid", "com.test.app")
	if err == nil {
		t.Error("expected error on non-darwin platform")
	}
}

func TestCodeSign_NonDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("test only runs on non-darwin platforms")
	}
	_, err := CodeSign(context.Background(), "/tmp/app.app", "", "", false)
	if err == nil {
		t.Error("expected error on non-darwin platform")
	}
}

func TestCodeSignVerify_NonDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("test only runs on non-darwin platforms")
	}
	_, err := CodeSignVerify(context.Background(), "/tmp/app.app")
	if err == nil {
		t.Error("expected error on non-darwin platform")
	}
}

func TestFindIdentity_NonDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("test only runs on non-darwin platforms")
	}
	_, err := FindIdentity(context.Background(), "")
	if err == nil {
		t.Error("expected error on non-darwin platform")
	}
}
