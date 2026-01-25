package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ReviewSubmissionState represents the state of a review submission.
type ReviewSubmissionState string

const (
	ReviewSubmissionStateReadyForReview    ReviewSubmissionState = "READY_FOR_REVIEW"
	ReviewSubmissionStateWaitingForReview  ReviewSubmissionState = "WAITING_FOR_REVIEW"
	ReviewSubmissionStateInReview          ReviewSubmissionState = "IN_REVIEW"
	ReviewSubmissionStateUnresolvedIssues  ReviewSubmissionState = "UNRESOLVED_ISSUES"
	ReviewSubmissionStateCanceling         ReviewSubmissionState = "CANCELING"
	ReviewSubmissionStateComplete          ReviewSubmissionState = "COMPLETE"
)

// ReviewSubmissionAttributes describes a review submission.
type ReviewSubmissionAttributes struct {
	Platform       Platform              `json:"platform,omitempty"`
	State          ReviewSubmissionState `json:"state,omitempty"`
	SubmittedDate  *string               `json:"submittedDate,omitempty"`
	Submitted      bool                  `json:"submitted,omitempty"`
}

// ReviewSubmissionResource represents a review submission resource.
type ReviewSubmissionResource struct {
	Type       ResourceType               `json:"type"`
	ID         string                     `json:"id"`
	Attributes ReviewSubmissionAttributes `json:"attributes"`
}

// ReviewSubmissionResponse is the response from review submission endpoints.
type ReviewSubmissionResponse struct {
	Data  ReviewSubmissionResource `json:"data"`
	Links Links                    `json:"links,omitempty"`
}

// ReviewSubmissionsResponse is the response from review submissions list.
type ReviewSubmissionsResponse struct {
	Data  []ReviewSubmissionResource `json:"data"`
	Links Links                      `json:"links,omitempty"`
}

// ReviewSubmissionCreateAttributes describes attributes for creating a review submission.
type ReviewSubmissionCreateAttributes struct {
	Platform Platform `json:"platform"`
}

// ReviewSubmissionCreateRelationships describes relationships for creating a review submission.
type ReviewSubmissionCreateRelationships struct {
	App *Relationship `json:"app"`
}

// ReviewSubmissionCreateData is the data for creating a review submission.
type ReviewSubmissionCreateData struct {
	Type          ResourceType                        `json:"type"`
	Attributes    ReviewSubmissionCreateAttributes    `json:"attributes"`
	Relationships ReviewSubmissionCreateRelationships `json:"relationships"`
}

// ReviewSubmissionCreateRequest is a request to create a review submission.
type ReviewSubmissionCreateRequest struct {
	Data ReviewSubmissionCreateData `json:"data"`
}

// ReviewSubmissionUpdateAttributes describes attributes for updating a review submission.
type ReviewSubmissionUpdateAttributes struct {
	Submitted  *bool    `json:"submitted,omitempty"`
	Canceled   *bool    `json:"canceled,omitempty"`
}

// ReviewSubmissionUpdateData is the data for updating a review submission.
type ReviewSubmissionUpdateData struct {
	Type       ResourceType                     `json:"type"`
	ID         string                           `json:"id"`
	Attributes ReviewSubmissionUpdateAttributes `json:"attributes"`
}

// ReviewSubmissionUpdateRequest is a request to update a review submission.
type ReviewSubmissionUpdateRequest struct {
	Data ReviewSubmissionUpdateData `json:"data"`
}

// ReviewSubmissionItemAttributes describes a review submission item.
type ReviewSubmissionItemAttributes struct {
	State          string `json:"state,omitempty"`
	ResolvedIssues bool   `json:"resolvedIssues,omitempty"`
	Removed        bool   `json:"removed,omitempty"`
}

// ReviewSubmissionItemResource represents a review submission item.
type ReviewSubmissionItemResource struct {
	Type       ResourceType                   `json:"type"`
	ID         string                         `json:"id"`
	Attributes ReviewSubmissionItemAttributes `json:"attributes"`
}

// ReviewSubmissionItemResponse is the response from review submission item endpoints.
type ReviewSubmissionItemResponse struct {
	Data  ReviewSubmissionItemResource `json:"data"`
	Links Links                        `json:"links,omitempty"`
}

// ReviewSubmissionItemCreateRelationships describes relationships for creating a submission item.
type ReviewSubmissionItemCreateRelationships struct {
	ReviewSubmission   *Relationship `json:"reviewSubmission"`
	AppStoreVersion    *Relationship `json:"appStoreVersion,omitempty"`
}

// ReviewSubmissionItemCreateData is the data for creating a review submission item.
type ReviewSubmissionItemCreateData struct {
	Type          ResourceType                            `json:"type"`
	Relationships ReviewSubmissionItemCreateRelationships `json:"relationships"`
}

// ReviewSubmissionItemCreateRequest is a request to create a review submission item.
type ReviewSubmissionItemCreateRequest struct {
	Data ReviewSubmissionItemCreateData `json:"data"`
}

// CreateReviewSubmission creates a new review submission for an app.
func (c *Client) CreateReviewSubmission(ctx context.Context, appID string, platform Platform) (*ReviewSubmissionResponse, error) {
	payload := ReviewSubmissionCreateRequest{
		Data: ReviewSubmissionCreateData{
			Type: ResourceTypeReviewSubmissions,
			Attributes: ReviewSubmissionCreateAttributes{
				Platform: platform,
			},
			Relationships: ReviewSubmissionCreateRelationships{
				App: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeApps,
						ID:   strings.TrimSpace(appID),
					},
				},
			},
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/reviewSubmissions", body)
	if err != nil {
		return nil, err
	}

	var response ReviewSubmissionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetReviewSubmission retrieves a review submission by ID.
func (c *Client) GetReviewSubmission(ctx context.Context, submissionID string) (*ReviewSubmissionResponse, error) {
	path := fmt.Sprintf("/v1/reviewSubmissions/%s", strings.TrimSpace(submissionID))
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response ReviewSubmissionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// SubmitReviewSubmission submits a review submission for review.
func (c *Client) SubmitReviewSubmission(ctx context.Context, submissionID string) (*ReviewSubmissionResponse, error) {
	submitted := true
	payload := ReviewSubmissionUpdateRequest{
		Data: ReviewSubmissionUpdateData{
			Type:       ResourceTypeReviewSubmissions,
			ID:         strings.TrimSpace(submissionID),
			Attributes: ReviewSubmissionUpdateAttributes{
				Submitted: &submitted,
			},
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/reviewSubmissions/%s", strings.TrimSpace(submissionID))
	data, err := c.do(ctx, "PATCH", path, body)
	if err != nil {
		return nil, err
	}

	var response ReviewSubmissionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CancelReviewSubmission cancels a review submission.
func (c *Client) CancelReviewSubmission(ctx context.Context, submissionID string) (*ReviewSubmissionResponse, error) {
	canceled := true
	payload := ReviewSubmissionUpdateRequest{
		Data: ReviewSubmissionUpdateData{
			Type:       ResourceTypeReviewSubmissions,
			ID:         strings.TrimSpace(submissionID),
			Attributes: ReviewSubmissionUpdateAttributes{
				Canceled: &canceled,
			},
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/reviewSubmissions/%s", strings.TrimSpace(submissionID))
	data, err := c.do(ctx, "PATCH", path, body)
	if err != nil {
		return nil, err
	}

	var response ReviewSubmissionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// AddReviewSubmissionItem adds an app store version to a review submission.
func (c *Client) AddReviewSubmissionItem(ctx context.Context, submissionID, versionID string) (*ReviewSubmissionItemResponse, error) {
	payload := ReviewSubmissionItemCreateRequest{
		Data: ReviewSubmissionItemCreateData{
			Type: ResourceTypeReviewSubmissionItems,
			Relationships: ReviewSubmissionItemCreateRelationships{
				ReviewSubmission: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeReviewSubmissions,
						ID:   strings.TrimSpace(submissionID),
					},
				},
				AppStoreVersion: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeAppStoreVersions,
						ID:   strings.TrimSpace(versionID),
					},
				},
			},
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/reviewSubmissionItems", body)
	if err != nil {
		return nil, err
	}

	var response ReviewSubmissionItemResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// DeleteReviewSubmissionItem removes an item from a review submission.
func (c *Client) DeleteReviewSubmissionItem(ctx context.Context, itemID string) error {
	path := fmt.Sprintf("/v1/reviewSubmissionItems/%s", strings.TrimSpace(itemID))
	_, err := c.do(ctx, "DELETE", path, nil)
	return err
}
