package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// AppStoreVersionAgeRatingDeclarationLinkageResponse is the response for age rating relationships.
type AppStoreVersionAgeRatingDeclarationLinkageResponse struct {
	Data  ResourceData `json:"data"`
	Links Links        `json:"links"`
}

// AppStoreVersionReviewDetailLinkageResponse is the response for review detail relationships.
type AppStoreVersionReviewDetailLinkageResponse struct {
	Data  ResourceData `json:"data"`
	Links Links        `json:"links"`
}

// AppStoreVersionAppClipDefaultExperienceLinkageResponse is the response for app clip default experience relationships.
type AppStoreVersionAppClipDefaultExperienceLinkageResponse struct {
	Data  ResourceData `json:"data"`
	Links Links        `json:"links"`
}

// AppStoreVersionSubmissionLinkageResponse is the response for submission relationships.
type AppStoreVersionSubmissionLinkageResponse struct {
	Data  ResourceData `json:"data"`
	Links Links        `json:"links"`
}

// AppStoreVersionRoutingAppCoverageLinkageResponse is the response for routing coverage relationships.
type AppStoreVersionRoutingAppCoverageLinkageResponse struct {
	Data  ResourceData `json:"data"`
	Links Links        `json:"links"`
}

// AppStoreVersionGameCenterAppVersionLinkageResponse is the response for Game Center app version relationships.
type AppStoreVersionGameCenterAppVersionLinkageResponse struct {
	Data  ResourceData `json:"data"`
	Links Links        `json:"links"`
}

// GetAppStoreVersionAgeRatingDeclarationRelationship retrieves the age rating linkage for a version.
func (c *Client) GetAppStoreVersionAgeRatingDeclarationRelationship(ctx context.Context, versionID string) (*AppStoreVersionAgeRatingDeclarationLinkageResponse, error) {
	versionID = strings.TrimSpace(versionID)
	if versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/ageRatingDeclaration", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionAgeRatingDeclarationLinkageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionReviewDetailRelationship retrieves the review detail linkage for a version.
func (c *Client) GetAppStoreVersionReviewDetailRelationship(ctx context.Context, versionID string) (*AppStoreVersionReviewDetailLinkageResponse, error) {
	versionID = strings.TrimSpace(versionID)
	if versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/appStoreReviewDetail", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionReviewDetailLinkageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionAppClipDefaultExperienceRelationship retrieves the app clip default experience linkage.
func (c *Client) GetAppStoreVersionAppClipDefaultExperienceRelationship(ctx context.Context, versionID string) (*AppStoreVersionAppClipDefaultExperienceLinkageResponse, error) {
	versionID = strings.TrimSpace(versionID)
	if versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/appClipDefaultExperience", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionAppClipDefaultExperienceLinkageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionExperimentsRelationships retrieves experiment linkages for a version (v1).
func (c *Client) GetAppStoreVersionExperimentsRelationships(ctx context.Context, versionID string, opts ...LinkagesOption) (*LinkagesResponse, error) {
	query := &linkagesQuery{}
	for _, opt := range opts {
		opt(query)
	}

	versionID = strings.TrimSpace(versionID)
	if query.nextURL == "" && versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/appStoreVersionExperiments", versionID)
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("appStoreVersionExperimentsRelationships: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildLinkagesQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response LinkagesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionExperimentsV2Relationships retrieves experiment linkages for a version (v2).
func (c *Client) GetAppStoreVersionExperimentsV2Relationships(ctx context.Context, versionID string, opts ...LinkagesOption) (*LinkagesResponse, error) {
	query := &linkagesQuery{}
	for _, opt := range opts {
		opt(query)
	}

	versionID = strings.TrimSpace(versionID)
	if query.nextURL == "" && versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/appStoreVersionExperimentsV2", versionID)
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("appStoreVersionExperimentsV2Relationships: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildLinkagesQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response LinkagesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionSubmissionRelationship retrieves the submission linkage for a version.
func (c *Client) GetAppStoreVersionSubmissionRelationship(ctx context.Context, versionID string) (*AppStoreVersionSubmissionLinkageResponse, error) {
	versionID = strings.TrimSpace(versionID)
	if versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/appStoreVersionSubmission", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionSubmissionLinkageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionCustomerReviewsRelationships retrieves customer review linkages for a version.
func (c *Client) GetAppStoreVersionCustomerReviewsRelationships(ctx context.Context, versionID string, opts ...LinkagesOption) (*LinkagesResponse, error) {
	query := &linkagesQuery{}
	for _, opt := range opts {
		opt(query)
	}

	versionID = strings.TrimSpace(versionID)
	if query.nextURL == "" && versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/customerReviews", versionID)
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("customerReviewsRelationships: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildLinkagesQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response LinkagesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionRoutingAppCoverageRelationship retrieves routing coverage linkage for a version.
func (c *Client) GetAppStoreVersionRoutingAppCoverageRelationship(ctx context.Context, versionID string) (*AppStoreVersionRoutingAppCoverageLinkageResponse, error) {
	versionID = strings.TrimSpace(versionID)
	if versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/routingAppCoverage", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionRoutingAppCoverageLinkageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionGameCenterAppVersionRelationship retrieves Game Center app version linkage.
func (c *Client) GetAppStoreVersionGameCenterAppVersionRelationship(ctx context.Context, versionID string) (*AppStoreVersionGameCenterAppVersionLinkageResponse, error) {
	versionID = strings.TrimSpace(versionID)
	if versionID == "" {
		return nil, fmt.Errorf("versionID is required")
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/gameCenterAppVersion", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionGameCenterAppVersionLinkageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
