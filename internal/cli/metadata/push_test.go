package metadata

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLocalMetadataTreatsDefaultLocaleCaseInsensitively(t *testing.T) {
	dir := t.TempDir()
	version := "1.2.3"

	if err := os.MkdirAll(filepath.Join(dir, appInfoDirName), 0o755); err != nil {
		t.Fatalf("mkdir app-info: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(dir, versionDirName, version), 0o755); err != nil {
		t.Fatalf("mkdir version dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, appInfoDirName, "Default.json"), []byte(`{"name":"Default App Name"}`), 0o644); err != nil {
		t.Fatalf("write app-info default file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, versionDirName, version, "DeFaUlT.json"), []byte(`{"description":"Default description"}`), 0o644); err != nil {
		t.Fatalf("write version default file: %v", err)
	}

	bundle, err := loadLocalMetadata(dir, version)
	if err != nil {
		t.Fatalf("loadLocalMetadata() error: %v", err)
	}
	if bundle.defaultAppInfo == nil {
		t.Fatal("expected default app-info localization")
	}
	if bundle.defaultVersion == nil {
		t.Fatal("expected default version localization")
	}
	if bundle.defaultAppInfo.localization.Name != "Default App Name" {
		t.Fatalf("expected default app-info name, got %q", bundle.defaultAppInfo.localization.Name)
	}
	if bundle.defaultVersion.localization.Description != "Default description" {
		t.Fatalf("expected default version description, got %q", bundle.defaultVersion.localization.Description)
	}
	if len(bundle.appInfo) != 0 {
		t.Fatalf("expected no explicit app-info locales, got %+v", bundle.appInfo)
	}
	if len(bundle.version) != 0 {
		t.Fatalf("expected no explicit version locales, got %+v", bundle.version)
	}
}

func TestLoadLocalMetadataRejectsVersionPathTraversal(t *testing.T) {
	dir := t.TempDir()

	_, err := loadLocalMetadata(dir, "../../secret")
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected usage error for invalid version, got %v", err)
	}
}

func TestBuildScopePlanCountsDeleteAndCreateForRecreate(t *testing.T) {
	local := map[string]localPlanFields{
		"en-US": {
			setFields: map[string]string{
				"name": "Local Name",
			},
			clearFields: map[string]struct{}{
				"subtitle": {},
			},
		},
	}
	remote := map[string]map[string]string{
		"en-US": {
			"name":     "Remote Name",
			"subtitle": "Remote subtitle",
		},
	}

	adds, updates, deletes, calls := buildScopePlan(
		appInfoDirName,
		"",
		appInfoPlanFields,
		local,
		remote,
	)

	if len(adds) != 0 {
		t.Fatalf("expected no adds, got %+v", adds)
	}
	if len(updates) != 1 {
		t.Fatalf("expected one field update, got %+v", updates)
	}
	if len(deletes) != 1 {
		t.Fatalf("expected one field delete, got %+v", deletes)
	}
	if calls.create != 1 || calls.delete != 1 || calls.update != 0 {
		t.Fatalf("unexpected call counts: %+v", calls)
	}
}

func TestApplyDefaultFallbackSkipsRemoteLocalesWhenDeletesAllowed(t *testing.T) {
	defaultAppInfo := appInfoLocalPatch{
		localization: AppInfoLocalization{Name: "Default Name"},
		setFields: map[string]string{
			"name": "Default Name",
		},
	}
	appInfo := applyDefaultAppInfoFallback(
		map[string]appInfoLocalPatch{},
		&defaultAppInfo,
		map[string]AppInfoLocalization{
			"fr": {Name: "Remote Name"},
		},
		true,
	)
	if len(appInfo) != 0 {
		t.Fatalf("expected no app-info fallback locales when deletes are allowed, got %+v", appInfo)
	}

	defaultVersion := versionLocalPatch{
		localization: VersionLocalization{Description: "Default Description"},
		setFields: map[string]string{
			"description": "Default Description",
		},
	}
	version := applyDefaultVersionFallback(
		map[string]versionLocalPatch{},
		&defaultVersion,
		map[string]VersionLocalization{
			"fr": {Description: "Remote Description"},
		},
		true,
	)
	if len(version) != 0 {
		t.Fatalf("expected no version fallback locales when deletes are allowed, got %+v", version)
	}
}

func TestReadAppInfoLocalizationPatchTracksClearsAndNoOps(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "en-US.json")
	if err := os.WriteFile(path, []byte(`{"name":"New Name","subtitle":"__ASC_DELETE__"}`), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	patch, err := readAppInfoLocalizationPatchFromFile(path)
	if err != nil {
		t.Fatalf("readAppInfoLocalizationPatchFromFile() error: %v", err)
	}
	if patch.localization.Name != "New Name" {
		t.Fatalf("expected name set, got %+v", patch.localization)
	}
	if _, ok := patch.clearFields["subtitle"]; !ok {
		t.Fatalf("expected subtitle clear marker, got %+v", patch.clearFields)
	}
}
