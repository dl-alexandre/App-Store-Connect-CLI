package asc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestGetGameCenterDetails_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails" {
			t.Fatalf("expected path /v1/gameCenterDetails, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "25" {
			t.Fatalf("expected limit=25, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetails(context.Background(), WithGCDetailsLimit(25)); err != nil {
		t.Fatalf("GetGameCenterDetails() error: %v", err)
	}
}

func TestGetGameCenterDetails_UsesNextURL(t *testing.T) {
	next := "https://api.appstoreconnect.apple.com/v1/gameCenterDetails?cursor=next"
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.String() != next {
			t.Fatalf("expected URL %q, got %q", next, req.URL.String())
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetails(context.Background(), WithGCDetailsNextURL(next)); err != nil {
		t.Fatalf("GetGameCenterDetails() error: %v", err)
	}
}

func TestGetGameCenterDetail(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":{"type":"gameCenterDetails","id":"detail-1","attributes":{"arcadeEnabled":true}}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails/detail-1" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1, got %s", req.URL.Path)
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetail(context.Background(), "detail-1"); err != nil {
		t.Fatalf("GetGameCenterDetail() error: %v", err)
	}
}

func TestGetGameCenterDetailGameCenterGroup(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":{"type":"gameCenterGroups","id":"group-1"}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/gameCenterGroup" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/gameCenterGroup, got %s", req.URL.Path)
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailGameCenterGroup(context.Background(), "detail-1"); err != nil {
		t.Fatalf("GetGameCenterDetailGameCenterGroup() error: %v", err)
	}
}

func TestGetGameCenterGroupGameCenterDetails_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterGroups/group-1/gameCenterDetails" {
			t.Fatalf("expected path /v1/gameCenterGroups/group-1/gameCenterDetails, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "30" {
			t.Fatalf("expected limit=30, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterGroupGameCenterDetails(context.Background(), "group-1", WithGCDetailsLimit(30)); err != nil {
		t.Fatalf("GetGameCenterGroupGameCenterDetails() error: %v", err)
	}
}

func TestGetGameCenterGroupGameCenterDetails_UsesNextURL(t *testing.T) {
	next := "https://api.appstoreconnect.apple.com/v1/gameCenterGroups/group-1/gameCenterDetails?cursor=next"
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.String() != next {
			t.Fatalf("expected URL %q, got %q", next, req.URL.String())
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterGroupGameCenterDetails(context.Background(), "", WithGCDetailsNextURL(next)); err != nil {
		t.Fatalf("GetGameCenterGroupGameCenterDetails() error: %v", err)
	}
}

func TestGCDetailsOptions(t *testing.T) {
	query := &gcDetailsQuery{}
	WithGCDetailsLimit(8)(query)
	if query.limit != 8 {
		t.Fatalf("expected limit 8, got %d", query.limit)
	}
	WithGCDetailsNextURL("next")(query)
	if query.nextURL != "next" {
		t.Fatalf("expected nextURL set, got %q", query.nextURL)
	}
	values, err := url.ParseQuery(buildGCDetailsQuery(query))
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if values.Get("limit") != "8" {
		t.Fatalf("expected limit=8, got %q", values.Get("limit"))
	}
}

func TestGetGameCenterDetailsAchievementReleases_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/achievementReleases" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/achievementReleases, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "12" {
			t.Fatalf("expected limit=12, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailsAchievementReleases(context.Background(), "detail-1", WithGCAchievementReleasesLimit(12)); err != nil {
		t.Fatalf("GetGameCenterDetailsAchievementReleases() error: %v", err)
	}
}

func TestGetGameCenterDetailsLeaderboardReleases_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/leaderboardReleases" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/leaderboardReleases, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "15" {
			t.Fatalf("expected limit=15, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailsLeaderboardReleases(context.Background(), "detail-1", WithGCLeaderboardReleasesLimit(15)); err != nil {
		t.Fatalf("GetGameCenterDetailsLeaderboardReleases() error: %v", err)
	}
}

func TestGetGameCenterDetailsLeaderboardSetReleases_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/leaderboardSetReleases" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/leaderboardSetReleases, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "18" {
			t.Fatalf("expected limit=18, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailsLeaderboardSetReleases(context.Background(), "detail-1", WithGCLeaderboardSetReleasesLimit(18)); err != nil {
		t.Fatalf("GetGameCenterDetailsLeaderboardSetReleases() error: %v", err)
	}
}

func TestGetGameCenterDetailsAchievementsV2_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/gameCenterAchievementsV2" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/gameCenterAchievementsV2, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "20" {
			t.Fatalf("expected limit=20, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailsAchievementsV2(context.Background(), "detail-1", WithGCAchievementsLimit(20)); err != nil {
		t.Fatalf("GetGameCenterDetailsAchievementsV2() error: %v", err)
	}
}

func TestGetGameCenterDetailsLeaderboardsV2_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/gameCenterLeaderboardsV2" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/gameCenterLeaderboardsV2, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "25" {
			t.Fatalf("expected limit=25, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailsLeaderboardsV2(context.Background(), "detail-1", WithGCLeaderboardsLimit(25)); err != nil {
		t.Fatalf("GetGameCenterDetailsLeaderboardsV2() error: %v", err)
	}
}

func TestGetGameCenterDetailsLeaderboardSetsV2_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/gameCenterLeaderboardSetsV2" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/gameCenterLeaderboardSetsV2, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "30" {
			t.Fatalf("expected limit=30, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailsLeaderboardSetsV2(context.Background(), "detail-1", WithGCLeaderboardSetsLimit(30)); err != nil {
		t.Fatalf("GetGameCenterDetailsLeaderboardSetsV2() error: %v", err)
	}
}

func TestGetGameCenterDetailsClassicMatchmakingRequests_WithQuery(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/metrics/classicMatchmakingRequests" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/metrics/classicMatchmakingRequests, got %s", req.URL.Path)
		}
		values := req.URL.Query()
		if values.Get("granularity") != "P1D" {
			t.Fatalf("expected granularity=P1D, got %q", values.Get("granularity"))
		}
		if values.Get("groupBy") != "result" {
			t.Fatalf("expected groupBy=result, got %q", values.Get("groupBy"))
		}
		if values.Get("filter[result]") != "MATCHED" {
			t.Fatalf("expected filter[result]=MATCHED, got %q", values.Get("filter[result]"))
		}
		if values.Get("sort") != "-count" {
			t.Fatalf("expected sort=-count, got %q", values.Get("sort"))
		}
		if values.Get("limit") != "50" {
			t.Fatalf("expected limit=50, got %q", values.Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	opts := []GCMatchmakingMetricsOption{
		WithGCMatchmakingMetricsGranularity("P1D"),
		WithGCMatchmakingMetricsGroupBy([]string{"result"}),
		WithGCMatchmakingMetricsFilterResult("MATCHED"),
		WithGCMatchmakingMetricsSort([]string{"-count"}),
		WithGCMatchmakingMetricsLimit(50),
	}

	if _, err := client.GetGameCenterDetailsClassicMatchmakingRequests(context.Background(), "detail-1", opts...); err != nil {
		t.Fatalf("GetGameCenterDetailsClassicMatchmakingRequests() error: %v", err)
	}
}

func TestGetGameCenterDetailsRuleBasedMatchmakingRequests_UsesNextURL(t *testing.T) {
	next := "https://api.appstoreconnect.apple.com/v1/gameCenterDetails/detail-1/metrics/ruleBasedMatchmakingRequests?cursor=next"
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.String() != next {
			t.Fatalf("expected URL %q, got %q", next, req.URL.String())
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailsRuleBasedMatchmakingRequests(context.Background(), "detail-1", WithGCMatchmakingMetricsNextURL(next)); err != nil {
		t.Fatalf("GetGameCenterDetailsRuleBasedMatchmakingRequests() error: %v", err)
	}
}

func TestCreateGameCenterDetail(t *testing.T) {
	response := jsonResponse(http.StatusCreated, `{"data":{"type":"gameCenterDetails","id":"detail-new","attributes":{"challengeEnabled":true}}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails" {
			t.Fatalf("expected path /v1/gameCenterDetails, got %s", req.URL.Path)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		var payload GameCenterDetailCreateRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		if payload.Data.Type != ResourceTypeGameCenterDetails {
			t.Fatalf("expected type gameCenterDetails, got %s", payload.Data.Type)
		}
		if payload.Data.Relationships == nil || payload.Data.Relationships.App == nil {
			t.Fatalf("expected app relationship to be set")
		}
		if payload.Data.Relationships.App.Data.ID != "app-123" {
			t.Fatalf("expected app ID app-123, got %s", payload.Data.Relationships.App.Data.ID)
		}
		if payload.Data.Relationships.App.Data.Type != ResourceTypeApps {
			t.Fatalf("expected app type apps, got %s", payload.Data.Relationships.App.Data.Type)
		}
		assertAuthorized(t, req)
	}, response)

	resp, err := client.CreateGameCenterDetail(context.Background(), "app-123", nil)
	if err != nil {
		t.Fatalf("CreateGameCenterDetail() error: %v", err)
	}

	if resp.Data.ID != "detail-new" {
		t.Fatalf("expected ID detail-new, got %s", resp.Data.ID)
	}
}

func TestUpdateGameCenterDetail(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":{"type":"gameCenterDetails","id":"detail-1","attributes":{"challengeEnabled":true}}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodPatch {
			t.Fatalf("expected PATCH, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails/detail-1" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1, got %s", req.URL.Path)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		var payload GameCenterDetailUpdateRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		if payload.Data.Type != ResourceTypeGameCenterDetails {
			t.Fatalf("expected type gameCenterDetails, got %s", payload.Data.Type)
		}
		if payload.Data.ID != "detail-1" {
			t.Fatalf("expected id detail-1, got %s", payload.Data.ID)
		}
		if payload.Data.Attributes == nil || payload.Data.Attributes.ChallengeEnabled == nil {
			t.Fatalf("expected challengeEnabled attribute to be set")
		}
		if *payload.Data.Attributes.ChallengeEnabled != true {
			t.Fatalf("expected challengeEnabled true, got %v", *payload.Data.Attributes.ChallengeEnabled)
		}
		assertAuthorized(t, req)
	}, response)

	enabled := true
	attrs := &GameCenterDetailUpdateAttributes{ChallengeEnabled: &enabled}
	resp, err := client.UpdateGameCenterDetail(context.Background(), "detail-1", attrs, nil)
	if err != nil {
		t.Fatalf("UpdateGameCenterDetail() error: %v", err)
	}

	if resp.Data.ID != "detail-1" {
		t.Fatalf("expected ID detail-1, got %s", resp.Data.ID)
	}
	if !resp.Data.Attributes.ChallengeEnabled {
		t.Fatalf("expected challengeEnabled true, got false")
	}
}
