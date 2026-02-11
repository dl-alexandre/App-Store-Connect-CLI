package cmdtest

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/validate"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/validation"
)

type validateFixture struct {
	app              string
	version          string
	appInfos         string
	appInfoLocs      string
	versionLocs      string
	ageRating        string
	screenshotSets   map[string]string
	screenshotsBySet map[string]string
}

func newValidateTestClient(t *testing.T, fixture validateFixture) *asc.Client {
	t.Helper()

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "key.p8")
	writeECDSAPEM(t, keyPath)

	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			return jsonResponse(http.StatusMethodNotAllowed, `{"errors":[{"status":405}]}`)
		}

		path := req.URL.Path
		switch {
		case path == "/v1/apps/app-1":
			return jsonResponse(http.StatusOK, fixture.app)
		case path == "/v1/appStoreVersions/ver-1":
			return jsonResponse(http.StatusOK, fixture.version)
		case path == "/v1/apps/app-1/appInfos":
			return jsonResponse(http.StatusOK, fixture.appInfos)
		case path == "/v1/appInfos/info-1/appInfoLocalizations":
			return jsonResponse(http.StatusOK, fixture.appInfoLocs)
		case path == "/v1/appStoreVersions/ver-1/appStoreVersionLocalizations":
			return jsonResponse(http.StatusOK, fixture.versionLocs)
		case path == "/v1/appStoreVersions/ver-1/ageRatingDeclaration":
			return jsonResponse(http.StatusOK, fixture.ageRating)
		case strings.HasPrefix(path, "/v1/appStoreVersionLocalizations/") && strings.HasSuffix(path, "/appScreenshotSets"):
			localizationID := strings.TrimSuffix(strings.TrimPrefix(path, "/v1/appStoreVersionLocalizations/"), "/appScreenshotSets")
			if body, ok := fixture.screenshotSets[localizationID]; ok {
				return jsonResponse(http.StatusOK, body)
			}
		case strings.HasPrefix(path, "/v1/appScreenshotSets/") && strings.HasSuffix(path, "/appScreenshots"):
			setID := strings.TrimSuffix(strings.TrimPrefix(path, "/v1/appScreenshotSets/"), "/appScreenshots")
			if body, ok := fixture.screenshotsBySet[setID]; ok {
				return jsonResponse(http.StatusOK, body)
			}
		}

		return jsonResponse(http.StatusNotFound, `{"errors":[{"status":404}]}`)
	})

	httpClient := &http.Client{Transport: transport}
	client, err := asc.NewClientWithHTTPClient("KEY123", "ISS456", keyPath, httpClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func jsonResponse(status int, body string) (*http.Response, error) {
	return &http.Response{
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func validValidateFixture() validateFixture {
	return validateFixture{
		app:         `{"data":{"type":"apps","id":"app-1","attributes":{"primaryLocale":"en-US"}}}`,
		version:     `{"data":{"type":"appStoreVersions","id":"ver-1","attributes":{"platform":"IOS","versionString":"1.0"}}}`,
		appInfos:    `{"data":[{"type":"appInfos","id":"info-1","attributes":{"state":"PREPARE_FOR_SUBMISSION"}}]}`,
		appInfoLocs: `{"data":[{"type":"appInfoLocalizations","id":"info-loc-1","attributes":{"locale":"en-US","name":"My App","subtitle":"Subtitle"}}]}`,
		versionLocs: `{"data":[{"type":"appStoreVersionLocalizations","id":"ver-loc-1","attributes":{"locale":"en-US","description":"Description","keywords":"keyword","whatsNew":"Notes","promotionalText":"Promo","supportUrl":"https://support.example.com","marketingUrl":"https://marketing.example.com"}}]}`,
		ageRating: `{"data":{"type":"ageRatingDeclarations","id":"age-1","attributes":{
			"advertising":false,
			"gambling":false,
			"healthOrWellnessTopics":false,
			"lootBox":false,
			"messagingAndChat":true,
			"parentalControls":true,
			"ageAssurance":false,
			"unrestrictedWebAccess":false,
			"userGeneratedContent":true,
			"alcoholTobaccoOrDrugUseOrReferences":"NONE",
			"contests":"NONE",
			"gamblingSimulated":"NONE",
			"gunsOrOtherWeapons":"NONE",
			"medicalOrTreatmentInformation":"NONE",
			"profanityOrCrudeHumor":"NONE",
			"sexualContentGraphicAndNudity":"NONE",
			"sexualContentOrNudity":"NONE",
			"horrorOrFearThemes":"NONE",
			"matureOrSuggestiveThemes":"NONE",
			"violenceCartoonOrFantasy":"NONE",
			"violenceRealistic":"NONE",
			"violenceRealisticProlongedGraphicOrSadistic":"NONE"
		}}}`,
		screenshotSets: map[string]string{
			"ver-loc-1": `{"data":[{"type":"appScreenshotSets","id":"set-1","attributes":{"screenshotDisplayType":"APP_IPHONE_65"}}]}`,
		},
		screenshotsBySet: map[string]string{
			"set-1": `{"data":[{"type":"appScreenshots","id":"shot-1","attributes":{"fileName":"shot.png","fileSize":1024,"imageAsset":{"width":1242,"height":2688}}}]}`,
		},
	}
}

func TestValidateRequiresAppAndVersionID(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing app",
			args:    []string{"validate", "--version-id", "ver-1"},
			wantErr: "--app is required",
		},
		{
			name:    "missing version id",
			args:    []string{"validate", "--app", "app-1"},
			wantErr: "--version-id is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			stdout, stderr := captureOutput(t, func() {
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
			if !strings.Contains(stderr, test.wantErr) {
				t.Fatalf("expected error %q, got %q", test.wantErr, stderr)
			}
		})
	}
}

func TestValidateOutputsJSONAndTable(t *testing.T) {
	fixture := validValidateFixture()
	client := newValidateTestClient(t, fixture)
	restore := validate.SetClientFactory(func() (*asc.Client, error) {
		return client, nil
	})
	defer restore()

	root := RootCommand("1.2.3")
	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"validate", "--app", "app-1", "--version-id", "ver-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	var report validation.Report
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}
	if report.Summary.Errors != 0 || report.Summary.Warnings != 0 {
		t.Fatalf("expected no issues, got %+v", report.Summary)
	}

	root = RootCommand("1.2.3")
	stdout, _ = captureOutput(t, func() {
		if err := root.Parse([]string{"validate", "--app", "app-1", "--version-id", "ver-1", "--output", "table"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	if !strings.Contains(stdout, "Severity") {
		t.Fatalf("expected table output to include headers, got %q", stdout)
	}
}

func TestValidateStrictExitBehavior(t *testing.T) {
	fixture := validValidateFixture()
	fixture.appInfoLocs = `{"data":[{"type":"appInfoLocalizations","id":"info-loc-1","attributes":{"locale":"en-US","name":"My App"}}]}`

	client := newValidateTestClient(t, fixture)
	restore := validate.SetClientFactory(func() (*asc.Client, error) {
		return client, nil
	})
	defer restore()

	root := RootCommand("1.2.3")
	_, _ = captureOutput(t, func() {
		if err := root.Parse([]string{"validate", "--app", "app-1", "--version-id", "ver-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	root = RootCommand("1.2.3")
	_, _ = captureOutput(t, func() {
		if err := root.Parse([]string{"validate", "--app", "app-1", "--version-id", "ver-1", "--strict"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if err == nil {
			t.Fatalf("expected error with --strict")
		}
		if _, ok := errors.AsType[ReportedError](err); !ok {
			t.Fatalf("expected ReportedError, got %v", err)
		}
	})
}

func TestValidateMixedWarningAndError(t *testing.T) {
	fixture := validValidateFixture()
	fixture.versionLocs = `{"data":[{"type":"appStoreVersionLocalizations","id":"ver-loc-1","attributes":{"locale":"en-US","description":"","keywords":"keyword","supportUrl":"https://support.example.com"}}]}`
	fixture.appInfoLocs = `{"data":[{"type":"appInfoLocalizations","id":"info-loc-1","attributes":{"locale":"en-US","name":"My App"}}]}`

	client := newValidateTestClient(t, fixture)
	restore := validate.SetClientFactory(func() (*asc.Client, error) {
		return client, nil
	})
	defer restore()

	root := RootCommand("1.2.3")
	_, _ = captureOutput(t, func() {
		if err := root.Parse([]string{"validate", "--app", "app-1", "--version-id", "ver-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if err == nil {
			t.Fatalf("expected error with mixed issues")
		}
		if _, ok := errors.AsType[ReportedError](err); !ok {
			t.Fatalf("expected ReportedError, got %v", err)
		}
	})
}
