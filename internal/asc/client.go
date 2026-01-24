package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

// FeedbackAttributes describes beta feedback screenshot submissions.
type FeedbackAttributes struct {
	CreatedDate    string                    `json:"createdDate"`
	Comment        string                    `json:"comment"`
	Email          string                    `json:"email"`
	DeviceModel    string                    `json:"deviceModel,omitempty"`
	OSVersion      string                    `json:"osVersion,omitempty"`
	AppPlatform    string                    `json:"appPlatform,omitempty"`
	DevicePlatform string                    `json:"devicePlatform,omitempty"`
	Screenshots    []FeedbackScreenshotImage `json:"screenshots,omitempty"`
}

// FeedbackScreenshotImage describes a screenshot attached to feedback.
type FeedbackScreenshotImage struct {
	URL            string `json:"url"`
	Width          int    `json:"width,omitempty"`
	Height         int    `json:"height,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
}

// CrashAttributes describes beta feedback crash submissions.
type CrashAttributes struct {
	CreatedDate    string `json:"createdDate"`
	Comment        string `json:"comment"`
	Email          string `json:"email"`
	DeviceModel    string `json:"deviceModel,omitempty"`
	OSVersion      string `json:"osVersion,omitempty"`
	AppPlatform    string `json:"appPlatform,omitempty"`
	DevicePlatform string `json:"devicePlatform,omitempty"`
	CrashLog       string `json:"crashLog,omitempty"`
}

// ReviewAttributes describes App Store customer reviews.
type ReviewAttributes struct {
	Rating           int    `json:"rating"`
	Title            string `json:"title"`
	Body             string `json:"body"`
	ReviewerNickname string `json:"reviewerNickname"`
	CreatedDate      string `json:"createdDate"`
	Territory        string `json:"territory"`
}

// FeedbackResponse is the response from beta feedback screenshots endpoint.
type FeedbackResponse = Response[FeedbackAttributes]

// CrashesResponse is the response from beta feedback crashes endpoint.
type CrashesResponse = Response[CrashAttributes]

// ReviewsResponse is the response from customer reviews endpoint.
type ReviewsResponse = Response[ReviewAttributes]

// AppsResponse is the response from apps endpoint.
type AppsResponse = Response[AppAttributes]

// AppResponse is the response from app detail endpoint.
type AppResponse = SingleResponse[AppAttributes]

// BuildsResponse is the response from builds endpoint.
type BuildsResponse = Response[BuildAttributes]

// AppStoreVersionsResponse is the response from app store versions endpoints.
type AppStoreVersionsResponse = Response[AppStoreVersionAttributes]

// AppStoreVersionResponse is the response from app store version detail.
type AppStoreVersionResponse = SingleResponse[AppStoreVersionAttributes]

// PreReleaseVersionsResponse is the response from pre-release versions endpoints.
type PreReleaseVersionsResponse struct {
	Data  []PreReleaseVersion `json:"data"`
	Links Links               `json:"links,omitempty"`
}

// PreReleaseVersionResponse is the response from pre-release version detail.
type PreReleaseVersionResponse struct {
	Data  PreReleaseVersion `json:"data"`
	Links Links             `json:"links,omitempty"`
}

// BuildResponse is the response from build detail/updates.
type BuildResponse = SingleResponse[BuildAttributes]

// AppStoreVersionLocalizationsResponse is the response from app store version localizations endpoints.
type AppStoreVersionLocalizationsResponse = Response[AppStoreVersionLocalizationAttributes]

// AppStoreVersionLocalizationResponse is the response from app store version localization detail/creates.
type AppStoreVersionLocalizationResponse = SingleResponse[AppStoreVersionLocalizationAttributes]

// AppInfoLocalizationsResponse is the response from app info localizations endpoints.
type AppInfoLocalizationsResponse = Response[AppInfoLocalizationAttributes]

// AppInfoLocalizationResponse is the response from app info localization detail/creates.
type AppInfoLocalizationResponse = SingleResponse[AppInfoLocalizationAttributes]

// AppInfosResponse is the response from app info endpoints.
type AppInfosResponse = Response[AppInfoAttributes]

// BetaGroupsResponse is the response from beta groups endpoints.
type BetaGroupsResponse = Response[BetaGroupAttributes]

// BetaGroupResponse is the response from beta group detail/creates.
type BetaGroupResponse = SingleResponse[BetaGroupAttributes]

// BetaTestersResponse is the response from beta testers endpoints.
type BetaTestersResponse = Response[BetaTesterAttributes]

// BetaTesterResponse is the response from beta tester detail/creates.
type BetaTesterResponse = SingleResponse[BetaTesterAttributes]

// BetaTesterInvitationResponse is the response from beta tester invitations.
type BetaTesterInvitationResponse = SingleResponse[struct{}]

// AppAttributes describes an app resource.
type AppAttributes struct {
	Name          string `json:"name"`
	BundleID      string `json:"bundleId"`
	SKU           string `json:"sku"`
	PrimaryLocale string `json:"primaryLocale,omitempty"`
}

// BuildAttributes describes a build resource.
type BuildAttributes struct {
	Version                 string `json:"version"`
	UploadedDate            string `json:"uploadedDate"`
	ExpirationDate          string `json:"expirationDate,omitempty"`
	ProcessingState         string `json:"processingState,omitempty"`
	MinOSVersion            string `json:"minOsVersion,omitempty"`
	UsesNonExemptEncryption bool   `json:"usesNonExemptEncryption,omitempty"`
	Expired                 bool   `json:"expired,omitempty"`
}

// AppStoreVersionAttributes describes app store version metadata.
type AppStoreVersionAttributes struct {
	Platform        Platform `json:"platform,omitempty"`
	VersionString   string   `json:"versionString,omitempty"`
	AppStoreState   string   `json:"appStoreState,omitempty"`
	AppVersionState string   `json:"appVersionState,omitempty"`
	CreatedDate     string   `json:"createdDate,omitempty"`
}

// PreReleaseVersionAttributes describes TestFlight pre-release version metadata.
type PreReleaseVersionAttributes struct {
	Version  string   `json:"version,omitempty"`
	Platform Platform `json:"platform,omitempty"`
}

// PreReleaseVersion represents a pre-release version resource.
type PreReleaseVersion struct {
	Type       ResourceType                `json:"type"`
	ID         string                      `json:"id"`
	Attributes PreReleaseVersionAttributes `json:"attributes"`
}

// AppStoreVersionLocalizationAttributes describes app store version localization metadata.
type AppStoreVersionLocalizationAttributes struct {
	Locale          string `json:"locale,omitempty"`
	Description     string `json:"description,omitempty"`
	Keywords        string `json:"keywords,omitempty"`
	MarketingURL    string `json:"marketingUrl,omitempty"`
	PromotionalText string `json:"promotionalText,omitempty"`
	SupportURL      string `json:"supportUrl,omitempty"`
	WhatsNew        string `json:"whatsNew,omitempty"`
}

// AppInfoLocalizationAttributes describes app info localization metadata.
type AppInfoLocalizationAttributes struct {
	Locale            string `json:"locale,omitempty"`
	Name              string `json:"name,omitempty"`
	Subtitle          string `json:"subtitle,omitempty"`
	PrivacyPolicyURL  string `json:"privacyPolicyUrl,omitempty"`
	PrivacyChoicesURL string `json:"privacyChoicesUrl,omitempty"`
	PrivacyPolicyText string `json:"privacyPolicyText,omitempty"`
}

// AppInfoAttributes describes app info resources.
type AppInfoAttributes struct{}

// BetaGroupAttributes describes a beta group resource.
type BetaGroupAttributes struct {
	Name                   string `json:"name"`
	CreatedDate            string `json:"createdDate,omitempty"`
	IsInternalGroup        bool   `json:"isInternalGroup,omitempty"`
	HasAccessToAllBuilds   bool   `json:"hasAccessToAllBuilds,omitempty"`
	PublicLinkEnabled      bool   `json:"publicLinkEnabled,omitempty"`
	PublicLinkLimitEnabled bool   `json:"publicLinkLimitEnabled,omitempty"`
	PublicLinkLimit        int    `json:"publicLinkLimit,omitempty"`
	PublicLink             string `json:"publicLink,omitempty"`
	FeedbackEnabled        bool   `json:"feedbackEnabled,omitempty"`
}

// BetaTesterAttributes describes a beta tester resource.
type BetaTesterAttributes struct {
	FirstName  string          `json:"firstName,omitempty"`
	LastName   string          `json:"lastName,omitempty"`
	Email      string          `json:"email,omitempty"`
	InviteType BetaInviteType  `json:"inviteType,omitempty"`
	State      BetaTesterState `json:"state,omitempty"`
}

// BetaInviteType represents the invitation type for a beta tester.
type BetaInviteType string

const (
	BetaInviteTypeEmail      BetaInviteType = "EMAIL"
	BetaInviteTypePublicLink BetaInviteType = "PUBLIC_LINK"
)

// BetaTesterState represents the invitation state for a beta tester.
type BetaTesterState string

const (
	BetaTesterStateNotInvited BetaTesterState = "NOT_INVITED"
	BetaTesterStateInvited    BetaTesterState = "INVITED"
	BetaTesterStateAccepted   BetaTesterState = "ACCEPTED"
	BetaTesterStateInstalled  BetaTesterState = "INSTALLED"
	BetaTesterStateRevoked    BetaTesterState = "REVOKED"
)

// BuildUploadAttributes describes a build upload resource.
type BuildUploadAttributes struct {
	CFBundleShortVersionString string   `json:"cfBundleShortVersionString"`
	CFBundleVersion            string   `json:"cfBundleVersion"`
	Platform                   Platform `json:"platform"`
	CreatedDate                *string  `json:"createdDate,omitempty"`
	State                      *string  `json:"state,omitempty"`
}

// BuildUploadRelationships describes the relationships for a build upload.
type BuildUploadRelationships struct {
	App   *Relationship `json:"app,omitempty"`
	Build *Relationship `json:"build,omitempty"`
}

// BuildUploadCreateData is the data portion of a build upload create request.
type BuildUploadCreateData struct {
	Type          ResourceType              `json:"type"`
	Attributes    BuildUploadAttributes     `json:"attributes"`
	Relationships *BuildUploadRelationships `json:"relationships,omitempty"`
}

// BuildUploadCreateRequest is a request to create a build upload.
type BuildUploadCreateRequest struct {
	Data BuildUploadCreateData `json:"data"`
}

// BuildUploadResponse is the response from build upload endpoint.
type BuildUploadResponse = SingleResourceResponse[BuildUploadAttributes]

// BuildUploadFileAttributes describes a build upload file resource.
type BuildUploadFileAttributes struct {
	AssetDeliveryState  *AppMediaAssetState `json:"assetDeliveryState,omitempty"`
	AssetToken          *string             `json:"assetToken,omitempty"`
	AssetType           AssetType           `json:"assetType,omitempty"`
	FileName            string              `json:"fileName"`
	FileSize            int64               `json:"fileSize"`
	SourceFileChecksums *Checksums          `json:"sourceFileChecksums,omitempty"`
	UploadOperations    []UploadOperation   `json:"uploadOperations,omitempty"`
	UTI                 UTI                 `json:"uti"`
	Uploaded            *bool               `json:"uploaded,omitempty"`
}

// BuildUploadFileRelationships describes the relationships for a build upload file.
type BuildUploadFileRelationships struct {
	BuildUpload *Relationship `json:"buildUpload"`
}

// BuildUploadFileCreateData is the data portion of a build upload file create request.
type BuildUploadFileCreateData struct {
	Type          ResourceType                  `json:"type"`
	Attributes    BuildUploadFileAttributes     `json:"attributes"`
	Relationships *BuildUploadFileRelationships `json:"relationships"`
}

// BuildUploadFileCreateRequest is a request to create a build upload file.
type BuildUploadFileCreateRequest struct {
	Data BuildUploadFileCreateData `json:"data"`
}

// BuildUploadFileResponse is the response from build upload file endpoint.
type BuildUploadFileResponse = SingleResourceResponse[BuildUploadFileAttributes]

// BuildUploadFileUpdateAttributes describes the attributes to update on a build upload file.
type BuildUploadFileUpdateAttributes struct {
	SourceFileChecksums *Checksums `json:"sourceFileChecksums,omitempty"`
	Uploaded            *bool      `json:"uploaded,omitempty"`
}

// BuildUploadFileUpdateData is the data portion of a build upload file update request.
type BuildUploadFileUpdateData struct {
	Type       ResourceType                     `json:"type"`
	ID         string                           `json:"id"`
	Attributes *BuildUploadFileUpdateAttributes `json:"attributes,omitempty"`
}

// BuildUploadFileUpdateRequest is a request to update a build upload file.
type BuildUploadFileUpdateRequest struct {
	Data BuildUploadFileUpdateData `json:"data"`
}

// UploadOperation represents a file upload operation with presigned URL.
type UploadOperation struct {
	Method         string       `json:"method"`
	URL            string       `json:"url"`
	Length         int64        `json:"length"`
	Offset         int64        `json:"offset"`
	RequestHeaders []HTTPHeader `json:"requestHeaders,omitempty"`
	Expiration     *string      `json:"expiration,omitempty"`
}

// HTTPHeader represents an HTTP header.
type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Checksums represents file checksums.
type Checksums struct {
	File      *Checksum `json:"file,omitempty"`
	Composite *Checksum `json:"composite,omitempty"`
}

// Checksum represents a single checksum.
type Checksum struct {
	Hash      string            `json:"hash"`
	Algorithm ChecksumAlgorithm `json:"algorithm"`
}

// AppMediaAssetState represents the state of an asset.
type AppMediaAssetState struct {
	State    *string       `json:"state,omitempty"`
	Errors   []StateDetail `json:"errors,omitempty"`
	Warnings []StateDetail `json:"warnings,omitempty"`
	Infos    []StateDetail `json:"infos,omitempty"`
}

// StateDetail represents details about a state (errors, warnings, infos).
type StateDetail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// AppStoreVersionSubmissionCreateData is the data portion of an app store version submission create request.
type AppStoreVersionSubmissionCreateData struct {
	Type          ResourceType                            `json:"type"`
	Relationships *AppStoreVersionSubmissionRelationships `json:"relationships"`
}

// AppStoreVersionSubmissionCreateRequest is a request to create an app store version submission.
type AppStoreVersionSubmissionCreateRequest struct {
	Data AppStoreVersionSubmissionCreateData `json:"data"`
}

// AppStoreVersionSubmissionRelationships describes the relationships for an app store version submission.
type AppStoreVersionSubmissionRelationships struct {
	AppStoreVersion *Relationship `json:"appStoreVersion"`
}

// AppStoreVersionSubmissionAttributes describes an app store version submission resource.
type AppStoreVersionSubmissionAttributes struct {
	CreatedDate *string `json:"createdDate,omitempty"`
}

// AppStoreVersionSubmissionResponse is the response from app store version submission endpoint.
type AppStoreVersionSubmissionResponse = SingleResourceResponse[AppStoreVersionSubmissionAttributes]

// AppStoreVersionSubmissionResource represents a submission with relationships.
type AppStoreVersionSubmissionResource struct {
	Type       ResourceType `json:"type"`
	ID         string       `json:"id"`
	Attributes struct {
		CreatedDate *string `json:"createdDate,omitempty"`
	} `json:"attributes,omitempty"`
	Relationships struct {
		AppStoreVersion *Relationship `json:"appStoreVersion,omitempty"`
	} `json:"relationships,omitempty"`
}

// AppStoreVersionSubmissionResourceResponse is a response containing a submission resource.
type AppStoreVersionSubmissionResourceResponse struct {
	Data AppStoreVersionSubmissionResource `json:"data"`
}

// AppStoreVersionBuildRelationshipUpdateRequest is a request to attach a build to a version.
type AppStoreVersionBuildRelationshipUpdateRequest struct {
	Data ResourceData `json:"data"`
}

// AppStoreVersionLocalizationCreateData is the data portion of a version localization create request.
type AppStoreVersionLocalizationCreateData struct {
	Type          ResourceType                              `json:"type"`
	Attributes    AppStoreVersionLocalizationAttributes     `json:"attributes"`
	Relationships *AppStoreVersionLocalizationRelationships `json:"relationships"`
}

// AppStoreVersionLocalizationCreateRequest is a request to create a version localization.
type AppStoreVersionLocalizationCreateRequest struct {
	Data AppStoreVersionLocalizationCreateData `json:"data"`
}

// AppStoreVersionLocalizationUpdateData is the data portion of a version localization update request.
type AppStoreVersionLocalizationUpdateData struct {
	Type       ResourceType                          `json:"type"`
	ID         string                                `json:"id"`
	Attributes AppStoreVersionLocalizationAttributes `json:"attributes"`
}

// AppStoreVersionLocalizationUpdateRequest is a request to update a version localization.
type AppStoreVersionLocalizationUpdateRequest struct {
	Data AppStoreVersionLocalizationUpdateData `json:"data"`
}

// AppStoreVersionLocalizationRelationships describes relationships for version localizations.
type AppStoreVersionLocalizationRelationships struct {
	AppStoreVersion *Relationship `json:"appStoreVersion"`
}

// AppInfoLocalizationCreateData is the data portion of an app info localization create request.
type AppInfoLocalizationCreateData struct {
	Type          ResourceType                      `json:"type"`
	Attributes    AppInfoLocalizationAttributes     `json:"attributes"`
	Relationships *AppInfoLocalizationRelationships `json:"relationships"`
}

// AppInfoLocalizationCreateRequest is a request to create an app info localization.
type AppInfoLocalizationCreateRequest struct {
	Data AppInfoLocalizationCreateData `json:"data"`
}

// AppInfoLocalizationUpdateData is the data portion of an app info localization update request.
type AppInfoLocalizationUpdateData struct {
	Type       ResourceType                  `json:"type"`
	ID         string                        `json:"id"`
	Attributes AppInfoLocalizationAttributes `json:"attributes"`
}

// AppInfoLocalizationUpdateRequest is a request to update an app info localization.
type AppInfoLocalizationUpdateRequest struct {
	Data AppInfoLocalizationUpdateData `json:"data"`
}

// AppInfoLocalizationRelationships describes relationships for app info localizations.
type AppInfoLocalizationRelationships struct {
	AppInfo *Relationship `json:"appInfo"`
}

// BetaGroupCreateData is the data portion of a beta group create request.
type BetaGroupCreateData struct {
	Type          ResourceType            `json:"type"`
	Attributes    BetaGroupAttributes     `json:"attributes"`
	Relationships *BetaGroupRelationships `json:"relationships"`
}

// BetaGroupCreateRequest is a request to create a beta group.
type BetaGroupCreateRequest struct {
	Data BetaGroupCreateData `json:"data"`
}

// BetaGroupUpdateAttributes describes attributes for updating a beta group.
type BetaGroupUpdateAttributes struct {
	Name                   string `json:"name,omitempty"`
	PublicLinkEnabled      *bool  `json:"publicLinkEnabled,omitempty"`
	PublicLinkLimitEnabled *bool  `json:"publicLinkLimitEnabled,omitempty"`
	PublicLinkLimit        int    `json:"publicLinkLimit,omitempty"`
	FeedbackEnabled        *bool  `json:"feedbackEnabled,omitempty"`
	IsInternalGroup        *bool  `json:"isInternalGroup,omitempty"`
	HasAccessToAllBuilds   *bool  `json:"hasAccessToAllBuilds,omitempty"`
}

// BetaGroupUpdateData is the data portion of a beta group update request.
type BetaGroupUpdateData struct {
	Type       ResourceType               `json:"type"`
	ID         string                     `json:"id"`
	Attributes *BetaGroupUpdateAttributes `json:"attributes,omitempty"`
}

// BetaGroupUpdateRequest is a request to update a beta group.
type BetaGroupUpdateRequest struct {
	Data BetaGroupUpdateData `json:"data"`
}

// BetaGroupRelationships describes relationships for beta groups.
type BetaGroupRelationships struct {
	App *Relationship `json:"app"`
}

// BetaTesterCreateAttributes describes attributes for creating a beta tester.
type BetaTesterCreateAttributes struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
}

// BetaTesterCreateRelationships describes relationships for beta tester creation.
type BetaTesterCreateRelationships struct {
	BetaGroups *RelationshipList `json:"betaGroups,omitempty"`
}

// BetaTesterCreateData is the data portion of a beta tester create request.
type BetaTesterCreateData struct {
	Type          ResourceType                   `json:"type"`
	Attributes    BetaTesterCreateAttributes     `json:"attributes"`
	Relationships *BetaTesterCreateRelationships `json:"relationships,omitempty"`
}

// BetaTesterCreateRequest is a request to create a beta tester.
type BetaTesterCreateRequest struct {
	Data BetaTesterCreateData `json:"data"`
}

// BetaTesterInvitationCreateRelationships describes relationships for invitations.
type BetaTesterInvitationCreateRelationships struct {
	App        *Relationship `json:"app"`
	BetaTester *Relationship `json:"betaTester,omitempty"`
}

// BetaTesterInvitationCreateData is the data portion of an invitation create request.
type BetaTesterInvitationCreateData struct {
	Type          ResourceType                             `json:"type"`
	Relationships *BetaTesterInvitationCreateRelationships `json:"relationships"`
}

// BetaTesterInvitationCreateRequest is a request to create a beta tester invitation.
type BetaTesterInvitationCreateRequest struct {
	Data BetaTesterInvitationCreateData `json:"data"`
}

// BuildUploadResult represents CLI output for build upload preparation.
type BuildUploadResult struct {
	UploadID   string            `json:"uploadId"`
	FileID     string            `json:"fileId"`
	FileName   string            `json:"fileName"`
	FileSize   int64             `json:"fileSize"`
	Operations []UploadOperation `json:"operations,omitempty"`
}

// AppStoreVersionSubmissionResult represents CLI output for submissions.
type AppStoreVersionSubmissionResult struct {
	SubmissionID string  `json:"submissionId"`
	CreatedDate  *string `json:"createdDate,omitempty"`
}

// AppStoreVersionSubmissionCreateResult represents CLI output for submission creation.
type AppStoreVersionSubmissionCreateResult struct {
	SubmissionID string  `json:"submissionId"`
	VersionID    string  `json:"versionId"`
	BuildID      string  `json:"buildId"`
	CreatedDate  *string `json:"createdDate,omitempty"`
}

// AppStoreVersionSubmissionStatusResult represents CLI output for submission status.
type AppStoreVersionSubmissionStatusResult struct {
	ID            string  `json:"id"`
	VersionID     string  `json:"versionId,omitempty"`
	VersionString string  `json:"versionString,omitempty"`
	Platform      string  `json:"platform,omitempty"`
	State         string  `json:"state,omitempty"`
	CreatedDate   *string `json:"createdDate,omitempty"`
}

// AppStoreVersionSubmissionCancelResult represents CLI output for submission cancellation.
type AppStoreVersionSubmissionCancelResult struct {
	ID        string `json:"id"`
	Cancelled bool   `json:"cancelled"`
}

// AppStoreVersionDetailResult represents CLI output for version details.
type AppStoreVersionDetailResult struct {
	ID            string `json:"id"`
	VersionString string `json:"versionString,omitempty"`
	Platform      string `json:"platform,omitempty"`
	State         string `json:"state,omitempty"`
	BuildID       string `json:"buildId,omitempty"`
	BuildVersion  string `json:"buildVersion,omitempty"`
	SubmissionID  string `json:"submissionId,omitempty"`
}

// AppStoreVersionAttachBuildResult represents CLI output for build attachment.
type AppStoreVersionAttachBuildResult struct {
	VersionID string `json:"versionId"`
	BuildID   string `json:"buildId"`
	Attached  bool   `json:"attached"`
}

// BuildBetaGroupsUpdateResult represents CLI output for build beta group updates.
type BuildBetaGroupsUpdateResult struct {
	BuildID  string   `json:"buildId"`
	GroupIDs []string `json:"groupIds"`
	Action   string   `json:"action"`
}

// BetaTesterInvitationResult represents CLI output for invitations.
type BetaTesterInvitationResult struct {
	InvitationID string `json:"invitationId"`
	TesterID     string `json:"testerId,omitempty"`
	AppID        string `json:"appId,omitempty"`
	Email        string `json:"email,omitempty"`
}

// BetaTesterDeleteResult represents CLI output for deletions.
type BetaTesterDeleteResult struct {
	ID      string `json:"id"`
	Email   string `json:"email,omitempty"`
	Deleted bool   `json:"deleted"`
}

// BetaTesterGroupsUpdateResult represents CLI output for beta tester group updates.
type BetaTesterGroupsUpdateResult struct {
	TesterID string   `json:"testerId"`
	GroupIDs []string `json:"groupIds"`
	Action   string   `json:"action"`
}

// AppStoreVersionLocalizationDeleteResult represents CLI output for localization deletions.
type AppStoreVersionLocalizationDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// LocalizationFileResult represents a localization file written or read.
type LocalizationFileResult struct {
	Locale string `json:"locale"`
	Path   string `json:"path"`
}

// LocalizationDownloadResult represents CLI output for localization downloads.
type LocalizationDownloadResult struct {
	Type       string                   `json:"type"`
	VersionID  string                   `json:"versionId,omitempty"`
	AppID      string                   `json:"appId,omitempty"`
	AppInfoID  string                   `json:"appInfoId,omitempty"`
	OutputPath string                   `json:"outputPath"`
	Files      []LocalizationFileResult `json:"files"`
}

// LocalizationUploadLocaleResult represents a per-locale upload result.
type LocalizationUploadLocaleResult struct {
	Locale         string `json:"locale"`
	Action         string `json:"action"`
	LocalizationID string `json:"localizationId,omitempty"`
}

// LocalizationUploadResult represents CLI output for localization uploads.
type LocalizationUploadResult struct {
	Type      string                           `json:"type"`
	VersionID string                           `json:"versionId,omitempty"`
	AppID     string                           `json:"appId,omitempty"`
	AppInfoID string                           `json:"appInfoId,omitempty"`
	DryRun    bool                             `json:"dryRun"`
	Results   []LocalizationUploadLocaleResult `json:"results"`
}

// GetFeedback retrieves TestFlight feedback
func (c *Client) GetFeedback(ctx context.Context, appID string, opts ...FeedbackOption) (*FeedbackResponse, error) {
	query := &feedbackQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/apps/%s/betaFeedbackScreenshotSubmissions", appID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("feedback: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildFeedbackQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response FeedbackResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetCrashes retrieves TestFlight crash reports
func (c *Client) GetCrashes(ctx context.Context, appID string, opts ...CrashOption) (*CrashesResponse, error) {
	query := &crashQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/apps/%s/betaFeedbackCrashSubmissions", appID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("crashes: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildCrashQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response CrashesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetReviews retrieves App Store reviews
func (c *Client) GetReviews(ctx context.Context, appID string, opts ...ReviewOption) (*ReviewsResponse, error) {
	query := &reviewQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/apps/%s/customerReviews", appID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("reviews: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildReviewQuery(opts); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response ReviewsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetApps retrieves the list of apps
func (c *Client) GetApps(ctx context.Context, opts ...AppsOption) (*AppsResponse, error) {
	query := &appsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := "/v1/apps"
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("apps: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildAppsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetApp retrieves a single app by ID.
func (c *Client) GetApp(ctx context.Context, appID string) (*AppResponse, error) {
	appID = strings.TrimSpace(appID)
	path := fmt.Sprintf("/v1/apps/%s", appID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBuilds retrieves the list of builds for an app
func (c *Client) GetBuilds(ctx context.Context, appID string, opts ...BuildsOption) (*BuildsResponse, error) {
	query := &buildsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/apps/%s/builds", appID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("builds: %w", err)
		}
		path = query.nextURL
	} else {
		values := url.Values{}
		// Use /v1/builds endpoint when sorting, limiting, or filtering by preReleaseVersion,
		// since /v1/apps/{id}/builds doesn't support these
		if query.sort != "" || query.limit > 0 || query.preReleaseVersionID != "" {
			path = "/v1/builds"
			values.Set("filter[app]", appID)
			if query.sort != "" {
				values.Set("sort", query.sort)
			}
			if query.limit > 0 {
				values.Set("limit", strconv.Itoa(query.limit))
			}
			if query.preReleaseVersionID != "" {
				values.Set("filter[preReleaseVersion]", query.preReleaseVersionID)
			}
		}
		if queryString := values.Encode(); queryString != "" {
			path += "?" + queryString
		}
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BuildsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersions retrieves app store versions for an app.
func (c *Client) GetAppStoreVersions(ctx context.Context, appID string, opts ...AppStoreVersionsOption) (*AppStoreVersionsResponse, error) {
	query := &appStoreVersionsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/apps/%s/appStoreVersions", appID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("appStoreVersions: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildAppStoreVersionsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetPreReleaseVersions retrieves TestFlight pre-release versions for an app.
func (c *Client) GetPreReleaseVersions(ctx context.Context, appID string, opts ...PreReleaseVersionsOption) (*PreReleaseVersionsResponse, error) {
	query := &preReleaseVersionsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	appID = strings.TrimSpace(appID)
	path := "/v1/preReleaseVersions"
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("preReleaseVersions: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildPreReleaseVersionsQuery(appID, query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response PreReleaseVersionsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetPreReleaseVersion retrieves a TestFlight pre-release version by ID.
func (c *Client) GetPreReleaseVersion(ctx context.Context, id string) (*PreReleaseVersionResponse, error) {
	id = strings.TrimSpace(id)
	path := fmt.Sprintf("/v1/preReleaseVersions/%s", id)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response PreReleaseVersionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersion retrieves an app store version by ID.
func (c *Client) GetAppStoreVersion(ctx context.Context, versionID string) (*AppStoreVersionResponse, error) {
	path := fmt.Sprintf("/v1/appStoreVersions/%s", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// AttachBuildToVersion attaches a build to an app store version.
func (c *Client) AttachBuildToVersion(ctx context.Context, versionID, buildID string) error {
	request := AppStoreVersionBuildRelationshipUpdateRequest{
		Data: ResourceData{
			Type: ResourceTypeBuilds,
			ID:   buildID,
		},
	}
	body, err := BuildRequestBody(request)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/relationships/build", versionID)
	if _, err := c.do(ctx, "PATCH", path, body); err != nil {
		return err
	}
	return nil
}

// GetAppStoreVersionBuild retrieves the build attached to a version.
func (c *Client) GetAppStoreVersionBuild(ctx context.Context, versionID string) (*BuildResponse, error) {
	path := fmt.Sprintf("/v1/appStoreVersions/%s/build", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BuildResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionSubmissionResource retrieves a submission by submission ID.
func (c *Client) GetAppStoreVersionSubmissionResource(ctx context.Context, submissionID string) (*AppStoreVersionSubmissionResourceResponse, error) {
	path := fmt.Sprintf("/v1/appStoreVersionSubmissions/%s", submissionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionSubmissionResourceResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionSubmissionForVersion retrieves a submission by version ID.
func (c *Client) GetAppStoreVersionSubmissionForVersion(ctx context.Context, versionID string) (*AppStoreVersionSubmissionResourceResponse, error) {
	path := fmt.Sprintf("/v1/appStoreVersions/%s/appStoreVersionSubmission", versionID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionSubmissionResourceResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBetaGroups retrieves the list of beta groups for an app.
func (c *Client) GetBetaGroups(ctx context.Context, appID string, opts ...BetaGroupsOption) (*BetaGroupsResponse, error) {
	query := &betaGroupsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/apps/%s/betaGroups", appID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("betaGroups: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildBetaGroupsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BetaGroupsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBetaGroupBuilds retrieves builds assigned to a beta group.
func (c *Client) GetBetaGroupBuilds(ctx context.Context, groupID string, opts ...BetaGroupBuildsOption) (*BuildsResponse, error) {
	query := &betaGroupBuildsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/betaGroups/%s/builds", groupID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("betaGroupBuilds: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildBetaGroupBuildsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BuildsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBetaGroupTesters retrieves beta testers assigned to a beta group.
func (c *Client) GetBetaGroupTesters(ctx context.Context, groupID string, opts ...BetaGroupTestersOption) (*BetaTestersResponse, error) {
	query := &betaGroupTestersQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/betaGroups/%s/betaTesters", groupID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("betaGroupTesters: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildBetaGroupTestersQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BetaTestersResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreateBetaGroup creates a beta group for an app.
func (c *Client) CreateBetaGroup(ctx context.Context, appID, name string) (*BetaGroupResponse, error) {
	payload := BetaGroupCreateRequest{
		Data: BetaGroupCreateData{
			Type:       ResourceTypeBetaGroups,
			Attributes: BetaGroupAttributes{Name: name},
			Relationships: &BetaGroupRelationships{
				App: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeApps,
						ID:   appID,
					},
				},
			},
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/betaGroups", body)
	if err != nil {
		return nil, err
	}

	var response BetaGroupResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBetaGroup retrieves a beta group by ID.
func (c *Client) GetBetaGroup(ctx context.Context, groupID string) (*BetaGroupResponse, error) {
	path := fmt.Sprintf("/v1/betaGroups/%s", groupID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BetaGroupResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// UpdateBetaGroup updates a beta group by ID.
func (c *Client) UpdateBetaGroup(ctx context.Context, groupID string, req BetaGroupUpdateRequest) (*BetaGroupResponse, error) {
	body, err := BuildRequestBody(req)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "PATCH", fmt.Sprintf("/v1/betaGroups/%s", groupID), body)
	if err != nil {
		return nil, err
	}

	var response BetaGroupResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// DeleteBetaGroup deletes a beta group by ID.
func (c *Client) DeleteBetaGroup(ctx context.Context, groupID string) error {
	path := fmt.Sprintf("/v1/betaGroups/%s", groupID)
	_, err := c.do(ctx, "DELETE", path, nil)
	return err
}

// AddBetaTestersToGroup adds testers to a beta group.
func (c *Client) AddBetaTestersToGroup(ctx context.Context, groupID string, testerIDs []string) error {
	testerIDs = normalizeList(testerIDs)
	payload := RelationshipRequest{
		Data: make([]RelationshipData, 0, len(testerIDs)),
	}
	for _, testerID := range testerIDs {
		payload.Data = append(payload.Data, RelationshipData{
			Type: ResourceTypeBetaTesters,
			ID:   testerID,
		})
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v1/betaGroups/%s/relationships/betaTesters", groupID)
	_, err = c.do(ctx, "POST", path, body)
	return err
}

// RemoveBetaTestersFromGroup removes testers from a beta group.
func (c *Client) RemoveBetaTestersFromGroup(ctx context.Context, groupID string, testerIDs []string) error {
	testerIDs = normalizeList(testerIDs)
	payload := RelationshipRequest{
		Data: make([]RelationshipData, 0, len(testerIDs)),
	}
	for _, testerID := range testerIDs {
		payload.Data = append(payload.Data, RelationshipData{
			Type: ResourceTypeBetaTesters,
			ID:   testerID,
		})
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v1/betaGroups/%s/relationships/betaTesters", groupID)
	_, err = c.do(ctx, "DELETE", path, body)
	return err
}

// GetBetaTesters retrieves beta testers for an app.
func (c *Client) GetBetaTesters(ctx context.Context, appID string, opts ...BetaTestersOption) (*BetaTestersResponse, error) {
	query := &betaTestersQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := "/v1/betaTesters"
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("betaTesters: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildBetaTestersQuery(appID, query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BetaTestersResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBetaTester retrieves a beta tester by ID.
func (c *Client) GetBetaTester(ctx context.Context, testerID string) (*BetaTesterResponse, error) {
	path := fmt.Sprintf("/v1/betaTesters/%s", testerID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BetaTesterResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreateBetaTester creates a beta tester.
func (c *Client) CreateBetaTester(ctx context.Context, email, firstName, lastName string, groupIDs []string) (*BetaTesterResponse, error) {
	groupIDs = normalizeList(groupIDs)
	var relationships *BetaTesterCreateRelationships
	if len(groupIDs) > 0 {
		relData := make([]ResourceData, 0, len(groupIDs))
		for _, groupID := range groupIDs {
			relData = append(relData, ResourceData{
				Type: ResourceTypeBetaGroups,
				ID:   groupID,
			})
		}
		relationships = &BetaTesterCreateRelationships{
			BetaGroups: &RelationshipList{Data: relData},
		}
	}

	payload := BetaTesterCreateRequest{
		Data: BetaTesterCreateData{
			Type: ResourceTypeBetaTesters,
			Attributes: BetaTesterCreateAttributes{
				FirstName: strings.TrimSpace(firstName),
				LastName:  strings.TrimSpace(lastName),
				Email:     strings.TrimSpace(email),
			},
			Relationships: relationships,
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/betaTesters", body)
	if err != nil {
		return nil, err
	}

	var response BetaTesterResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// AddBetaTesterToGroups adds a tester to multiple beta groups.
func (c *Client) AddBetaTesterToGroups(ctx context.Context, testerID string, groupIDs []string) error {
	testerID = strings.TrimSpace(testerID)
	groupIDs = normalizeList(groupIDs)
	if testerID == "" {
		return fmt.Errorf("tester ID is required")
	}
	if len(groupIDs) == 0 {
		return fmt.Errorf("group IDs are required")
	}

	relData := make([]ResourceData, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		relData = append(relData, ResourceData{
			Type: ResourceTypeBetaGroups,
			ID:   groupID,
		})
	}

	payload := RelationshipList{Data: relData}
	body, err := BuildRequestBody(payload)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v1/betaTesters/%s/relationships/betaGroups", testerID)
	if _, err := c.do(ctx, "POST", path, body); err != nil {
		return err
	}
	return nil
}

// RemoveBetaTesterFromGroups removes a tester from multiple beta groups.
func (c *Client) RemoveBetaTesterFromGroups(ctx context.Context, testerID string, groupIDs []string) error {
	testerID = strings.TrimSpace(testerID)
	groupIDs = normalizeList(groupIDs)
	if testerID == "" {
		return fmt.Errorf("tester ID is required")
	}
	if len(groupIDs) == 0 {
		return fmt.Errorf("group IDs are required")
	}

	relData := make([]ResourceData, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		relData = append(relData, ResourceData{
			Type: ResourceTypeBetaGroups,
			ID:   groupID,
		})
	}

	payload := RelationshipList{Data: relData}
	body, err := BuildRequestBody(payload)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v1/betaTesters/%s/relationships/betaGroups", testerID)
	if _, err := c.do(ctx, "DELETE", path, body); err != nil {
		return err
	}
	return nil
}

// DeleteBetaTester deletes a beta tester by ID.
func (c *Client) DeleteBetaTester(ctx context.Context, testerID string) error {
	path := fmt.Sprintf("/v1/betaTesters/%s", testerID)
	_, err := c.do(ctx, "DELETE", path, nil)
	return err
}

// CreateBetaTesterInvitation creates a beta tester invitation.
func (c *Client) CreateBetaTesterInvitation(ctx context.Context, appID, testerID string) (*BetaTesterInvitationResponse, error) {
	payload := BetaTesterInvitationCreateRequest{
		Data: BetaTesterInvitationCreateData{
			Type: ResourceTypeBetaTesterInvitations,
			Relationships: &BetaTesterInvitationCreateRelationships{
				App: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeApps,
						ID:   appID,
					},
				},
				BetaTester: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeBetaTesters,
						ID:   testerID,
					},
				},
			},
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/betaTesterInvitations", body)
	if err != nil {
		return nil, err
	}

	var response BetaTesterInvitationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionLocalizations retrieves localizations for an app store version.
func (c *Client) GetAppStoreVersionLocalizations(ctx context.Context, versionID string, opts ...AppStoreVersionLocalizationsOption) (*AppStoreVersionLocalizationsResponse, error) {
	query := &appStoreVersionLocalizationsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/appStoreVersions/%s/appStoreVersionLocalizations", versionID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("appStoreVersionLocalizations: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildAppStoreVersionLocalizationsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionLocalizationsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionLocalization retrieves a single app store version localization by ID.
func (c *Client) GetAppStoreVersionLocalization(ctx context.Context, localizationID string) (*AppStoreVersionLocalizationResponse, error) {
	path := fmt.Sprintf("/v1/appStoreVersionLocalizations/%s", localizationID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionLocalizationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreateAppStoreVersionLocalization creates a localization for an app store version.
func (c *Client) CreateAppStoreVersionLocalization(ctx context.Context, versionID string, attributes AppStoreVersionLocalizationAttributes) (*AppStoreVersionLocalizationResponse, error) {
	payload := AppStoreVersionLocalizationCreateRequest{
		Data: AppStoreVersionLocalizationCreateData{
			Type:       ResourceTypeAppStoreVersionLocalizations,
			Attributes: attributes,
			Relationships: &AppStoreVersionLocalizationRelationships{
				AppStoreVersion: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeAppStoreVersions,
						ID:   versionID,
					},
				},
			},
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/appStoreVersionLocalizations", body)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionLocalizationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// UpdateAppStoreVersionLocalization updates a localization for an app store version.
func (c *Client) UpdateAppStoreVersionLocalization(ctx context.Context, localizationID string, attributes AppStoreVersionLocalizationAttributes) (*AppStoreVersionLocalizationResponse, error) {
	payload := AppStoreVersionLocalizationUpdateRequest{
		Data: AppStoreVersionLocalizationUpdateData{
			Type:       ResourceTypeAppStoreVersionLocalizations,
			ID:         localizationID,
			Attributes: attributes,
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/appStoreVersionLocalizations/%s", localizationID)
	data, err := c.do(ctx, "PATCH", path, body)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionLocalizationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// DeleteAppStoreVersionLocalization deletes a localization by ID.
func (c *Client) DeleteAppStoreVersionLocalization(ctx context.Context, localizationID string) error {
	path := fmt.Sprintf("/v1/appStoreVersionLocalizations/%s", localizationID)
	if _, err := c.do(ctx, "DELETE", path, nil); err != nil {
		return err
	}
	return nil
}

// GetAppInfoLocalizations retrieves localizations for an app info resource.
func (c *Client) GetAppInfoLocalizations(ctx context.Context, appInfoID string, opts ...AppInfoLocalizationsOption) (*AppInfoLocalizationsResponse, error) {
	query := &appInfoLocalizationsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/appInfos/%s/appInfoLocalizations", appInfoID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("appInfoLocalizations: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildAppInfoLocalizationsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppInfoLocalizationsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreateAppInfoLocalization creates a localization for an app info resource.
func (c *Client) CreateAppInfoLocalization(ctx context.Context, appInfoID string, attributes AppInfoLocalizationAttributes) (*AppInfoLocalizationResponse, error) {
	payload := AppInfoLocalizationCreateRequest{
		Data: AppInfoLocalizationCreateData{
			Type:       ResourceTypeAppInfoLocalizations,
			Attributes: attributes,
			Relationships: &AppInfoLocalizationRelationships{
				AppInfo: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeAppInfos,
						ID:   appInfoID,
					},
				},
			},
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/appInfoLocalizations", body)
	if err != nil {
		return nil, err
	}

	var response AppInfoLocalizationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// UpdateAppInfoLocalization updates a localization for an app info resource.
func (c *Client) UpdateAppInfoLocalization(ctx context.Context, localizationID string, attributes AppInfoLocalizationAttributes) (*AppInfoLocalizationResponse, error) {
	payload := AppInfoLocalizationUpdateRequest{
		Data: AppInfoLocalizationUpdateData{
			Type:       ResourceTypeAppInfoLocalizations,
			ID:         localizationID,
			Attributes: attributes,
		},
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/appInfoLocalizations/%s", localizationID)
	data, err := c.do(ctx, "PATCH", path, body)
	if err != nil {
		return nil, err
	}

	var response AppInfoLocalizationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppInfos retrieves app info records for an app.
func (c *Client) GetAppInfos(ctx context.Context, appID string) (*AppInfosResponse, error) {
	path := fmt.Sprintf("/v1/apps/%s/appInfos", appID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppInfosResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBuild retrieves a single build by ID.
func (c *Client) GetBuild(ctx context.Context, buildID string) (*BuildResponse, error) {
	path := fmt.Sprintf("/v1/builds/%s", buildID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BuildResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBuildAppStoreVersion retrieves the app store version for a build.
func (c *Client) GetBuildAppStoreVersion(ctx context.Context, buildID string) (*AppStoreVersionResponse, error) {
	path := fmt.Sprintf("/v1/builds/%s/appStoreVersion", buildID)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// ExpireBuild expires a build for TestFlight testing.
func (c *Client) ExpireBuild(ctx context.Context, buildID string) (*BuildResponse, error) {
	payload := struct {
		Data struct {
			Type       ResourceType `json:"type"`
			ID         string       `json:"id"`
			Attributes struct {
				Expired bool `json:"expired"`
			} `json:"attributes"`
		} `json:"data"`
	}{}
	payload.Data.Type = ResourceTypeBuilds
	payload.Data.ID = buildID
	payload.Data.Attributes.Expired = true

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/builds/%s", buildID)
	data, err := c.do(ctx, "PATCH", path, body)
	if err != nil {
		return nil, err
	}

	var response BuildResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// AddBetaGroupsToBuild adds beta groups to a build for TestFlight distribution.
func (c *Client) AddBetaGroupsToBuild(ctx context.Context, buildID string, groupIDs []string) error {
	payload := RelationshipRequest{
		Data: make([]RelationshipData, len(groupIDs)),
	}
	for i, id := range groupIDs {
		payload.Data[i] = RelationshipData{
			Type: ResourceTypeBetaGroups,
			ID:   id,
		}
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v1/builds/%s/relationships/betaGroups", buildID)
	if _, err := c.do(ctx, "POST", path, body); err != nil {
		return err
	}
	return nil
}

// RemoveBetaGroupsFromBuild removes beta groups from a build.
func (c *Client) RemoveBetaGroupsFromBuild(ctx context.Context, buildID string, groupIDs []string) error {
	payload := RelationshipRequest{
		Data: make([]RelationshipData, len(groupIDs)),
	}
	for i, id := range groupIDs {
		payload.Data[i] = RelationshipData{
			Type: ResourceTypeBetaGroups,
			ID:   id,
		}
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v1/builds/%s/relationships/betaGroups", buildID)
	if _, err := c.do(ctx, "DELETE", path, body); err != nil {
		return err
	}
	return nil
}

// CreateBuildUpload creates a new build upload record.
func (c *Client) CreateBuildUpload(ctx context.Context, req BuildUploadCreateRequest) (*BuildUploadResponse, error) {
	body, err := BuildRequestBody(req)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/buildUploads", body)
	if err != nil {
		return nil, err
	}

	var response BuildUploadResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBuildUpload retrieves a build upload by ID.
func (c *Client) GetBuildUpload(ctx context.Context, id string) (*BuildUploadResponse, error) {
	data, err := c.do(ctx, "GET", fmt.Sprintf("/v1/buildUploads/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var response BuildUploadResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreateBuildUploadFile creates a new build upload file reservation.
func (c *Client) CreateBuildUploadFile(ctx context.Context, req BuildUploadFileCreateRequest) (*BuildUploadFileResponse, error) {
	body, err := BuildRequestBody(req)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/buildUploadFiles", body)
	if err != nil {
		return nil, err
	}

	var response BuildUploadFileResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// UpdateBuildUploadFile updates a build upload file (used to commit upload).
func (c *Client) UpdateBuildUploadFile(ctx context.Context, id string, req BuildUploadFileUpdateRequest) (*BuildUploadFileResponse, error) {
	body, err := BuildRequestBody(req)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "PATCH", fmt.Sprintf("/v1/buildUploadFiles/%s", id), body)
	if err != nil {
		return nil, err
	}

	var response BuildUploadFileResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreateAppStoreVersionSubmission creates a new app store version submission.
func (c *Client) CreateAppStoreVersionSubmission(ctx context.Context, req AppStoreVersionSubmissionCreateRequest) (*AppStoreVersionSubmissionResponse, error) {
	body, err := BuildRequestBody(req)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/appStoreVersionSubmissions", body)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionSubmissionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAppStoreVersionSubmission retrieves an app store version submission by ID.
func (c *Client) GetAppStoreVersionSubmission(ctx context.Context, id string) (*AppStoreVersionSubmissionResponse, error) {
	data, err := c.do(ctx, "GET", fmt.Sprintf("/v1/appStoreVersionSubmissions/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var response AppStoreVersionSubmissionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// DeleteAppStoreVersionSubmission deletes an app store version submission.
func (c *Client) DeleteAppStoreVersionSubmission(ctx context.Context, id string) error {
	_, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/appStoreVersionSubmissions/%s", id), nil)
	return err
}

// PrintJSON prints data as minified JSON (best for AI agents)
func PrintJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(data)
}

// PrintPrettyJSON prints data as indented JSON (best for debugging).
func PrintPrettyJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// PrintMarkdown prints data as Markdown table
func PrintMarkdown(data interface{}) error {
	switch v := data.(type) {
	case *FeedbackResponse:
		return printFeedbackMarkdown(v)
	case *CrashesResponse:
		return printCrashesMarkdown(v)
	case *ReviewsResponse:
		return printReviewsMarkdown(v)
	case *AppsResponse:
		return printAppsMarkdown(v)
	case *AppResponse:
		return printAppsMarkdown(&AppsResponse{Data: []Resource[AppAttributes]{v.Data}})
	case *BuildsResponse:
		return printBuildsMarkdown(v)
	case *AppStoreVersionsResponse:
		return printAppStoreVersionsMarkdown(v)
	case *PreReleaseVersionsResponse:
		return printPreReleaseVersionsMarkdown(v)
	case *BuildResponse:
		return printBuildsMarkdown(&BuildsResponse{Data: []Resource[BuildAttributes]{v.Data}})
	case *PreReleaseVersionResponse:
		return printPreReleaseVersionsMarkdown(&PreReleaseVersionsResponse{Data: []PreReleaseVersion{v.Data}})
	case *AppStoreVersionLocalizationsResponse:
		return printAppStoreVersionLocalizationsMarkdown(v)
	case *AppStoreVersionLocalizationResponse:
		return printAppStoreVersionLocalizationsMarkdown(&AppStoreVersionLocalizationsResponse{Data: []Resource[AppStoreVersionLocalizationAttributes]{v.Data}})
	case *AppInfoLocalizationsResponse:
		return printAppInfoLocalizationsMarkdown(v)
	case *BetaGroupsResponse:
		return printBetaGroupsMarkdown(v)
	case *BetaGroupResponse:
		return printBetaGroupsMarkdown(&BetaGroupsResponse{Data: []Resource[BetaGroupAttributes]{v.Data}})
	case *BetaTestersResponse:
		return printBetaTestersMarkdown(v)
	case *BetaTesterResponse:
		return printBetaTesterMarkdown(v)
	case *SandboxTestersResponse:
		return printSandboxTestersMarkdown(v)
	case *SandboxTesterResponse:
		return printSandboxTestersMarkdown(&SandboxTestersResponse{Data: []Resource[SandboxTesterAttributes]{v.Data}})
	case *LocalizationDownloadResult:
		return printLocalizationDownloadResultMarkdown(v)
	case *LocalizationUploadResult:
		return printLocalizationUploadResultMarkdown(v)
	case *BuildUploadResult:
		return printBuildUploadResultMarkdown(v)
	case *SalesReportResult:
		return printSalesReportResultMarkdown(v)
	case *FinanceReportResult:
		return printFinanceReportResultMarkdown(v)
	case *FinanceRegionsResult:
		return printFinanceRegionsMarkdown(v)
	case *AnalyticsReportRequestResult:
		return printAnalyticsReportRequestResultMarkdown(v)
	case *AnalyticsReportRequestsResponse:
		return printAnalyticsReportRequestsMarkdown(v)
	case *AnalyticsReportRequestResponse:
		return printAnalyticsReportRequestsMarkdown(&AnalyticsReportRequestsResponse{Data: []AnalyticsReportRequestResource{v.Data}, Links: v.Links})
	case *AnalyticsReportDownloadResult:
		return printAnalyticsReportDownloadResultMarkdown(v)
	case *AnalyticsReportGetResult:
		return printAnalyticsReportGetResultMarkdown(v)
	case *AppStoreVersionSubmissionResult:
		return printAppStoreVersionSubmissionMarkdown(v)
	case *AppStoreVersionSubmissionCreateResult:
		return printAppStoreVersionSubmissionCreateMarkdown(v)
	case *AppStoreVersionSubmissionStatusResult:
		return printAppStoreVersionSubmissionStatusMarkdown(v)
	case *AppStoreVersionSubmissionCancelResult:
		return printAppStoreVersionSubmissionCancelMarkdown(v)
	case *AppStoreVersionDetailResult:
		return printAppStoreVersionDetailMarkdown(v)
	case *AppStoreVersionAttachBuildResult:
		return printAppStoreVersionAttachBuildMarkdown(v)
	case *BuildBetaGroupsUpdateResult:
		return printBuildBetaGroupsUpdateMarkdown(v)
	case *BetaTesterDeleteResult:
		return printBetaTesterDeleteResultMarkdown(v)
	case *BetaTesterGroupsUpdateResult:
		return printBetaTesterGroupsUpdateResultMarkdown(v)
	case *AppStoreVersionLocalizationDeleteResult:
		return printAppStoreVersionLocalizationDeleteResultMarkdown(v)
	case *BetaTesterInvitationResult:
		return printBetaTesterInvitationResultMarkdown(v)
	case *SandboxTesterDeleteResult:
		return printSandboxTesterDeleteResultMarkdown(v)
	case *SandboxTesterClearHistoryResult:
		return printSandboxTesterClearHistoryResultMarkdown(v)
	case *XcodeCloudRunResult:
		return printXcodeCloudRunResultMarkdown(v)
	case *XcodeCloudStatusResult:
		return printXcodeCloudStatusResultMarkdown(v)
	case *CiProductsResponse:
		return printCiProductsMarkdown(v)
	case *CiWorkflowsResponse:
		return printCiWorkflowsMarkdown(v)
	case *CiBuildRunsResponse:
		return printCiBuildRunsMarkdown(v)
	case *CustomerReviewResponseResponse:
		return printCustomerReviewResponseMarkdown(v)
	case *CustomerReviewResponseDeleteResult:
		return printCustomerReviewResponseDeleteResultMarkdown(v)
	default:
		return PrintJSON(data)
	}
}

// PrintTable prints data as a formatted table
func PrintTable(data interface{}) error {
	switch v := data.(type) {
	case *FeedbackResponse:
		return printFeedbackTable(v)
	case *CrashesResponse:
		return printCrashesTable(v)
	case *ReviewsResponse:
		return printReviewsTable(v)
	case *AppsResponse:
		return printAppsTable(v)
	case *AppResponse:
		return printAppsTable(&AppsResponse{Data: []Resource[AppAttributes]{v.Data}})
	case *BuildsResponse:
		return printBuildsTable(v)
	case *AppStoreVersionsResponse:
		return printAppStoreVersionsTable(v)
	case *PreReleaseVersionsResponse:
		return printPreReleaseVersionsTable(v)
	case *BuildResponse:
		return printBuildsTable(&BuildsResponse{Data: []Resource[BuildAttributes]{v.Data}})
	case *PreReleaseVersionResponse:
		return printPreReleaseVersionsTable(&PreReleaseVersionsResponse{Data: []PreReleaseVersion{v.Data}})
	case *AppStoreVersionLocalizationsResponse:
		return printAppStoreVersionLocalizationsTable(v)
	case *AppStoreVersionLocalizationResponse:
		return printAppStoreVersionLocalizationsTable(&AppStoreVersionLocalizationsResponse{Data: []Resource[AppStoreVersionLocalizationAttributes]{v.Data}})
	case *AppInfoLocalizationsResponse:
		return printAppInfoLocalizationsTable(v)
	case *BetaGroupsResponse:
		return printBetaGroupsTable(v)
	case *BetaGroupResponse:
		return printBetaGroupsTable(&BetaGroupsResponse{Data: []Resource[BetaGroupAttributes]{v.Data}})
	case *BetaTestersResponse:
		return printBetaTestersTable(v)
	case *BetaTesterResponse:
		return printBetaTesterTable(v)
	case *SandboxTestersResponse:
		return printSandboxTestersTable(v)
	case *SandboxTesterResponse:
		return printSandboxTestersTable(&SandboxTestersResponse{Data: []Resource[SandboxTesterAttributes]{v.Data}})
	case *LocalizationDownloadResult:
		return printLocalizationDownloadResultTable(v)
	case *LocalizationUploadResult:
		return printLocalizationUploadResultTable(v)
	case *BuildUploadResult:
		return printBuildUploadResultTable(v)
	case *SalesReportResult:
		return printSalesReportResultTable(v)
	case *FinanceReportResult:
		return printFinanceReportResultTable(v)
	case *FinanceRegionsResult:
		return printFinanceRegionsTable(v)
	case *AnalyticsReportRequestResult:
		return printAnalyticsReportRequestResultTable(v)
	case *AnalyticsReportRequestsResponse:
		return printAnalyticsReportRequestsTable(v)
	case *AnalyticsReportRequestResponse:
		return printAnalyticsReportRequestsTable(&AnalyticsReportRequestsResponse{Data: []AnalyticsReportRequestResource{v.Data}, Links: v.Links})
	case *AnalyticsReportDownloadResult:
		return printAnalyticsReportDownloadResultTable(v)
	case *AnalyticsReportGetResult:
		return printAnalyticsReportGetResultTable(v)
	case *AppStoreVersionSubmissionResult:
		return printAppStoreVersionSubmissionTable(v)
	case *AppStoreVersionSubmissionCreateResult:
		return printAppStoreVersionSubmissionCreateTable(v)
	case *AppStoreVersionSubmissionStatusResult:
		return printAppStoreVersionSubmissionStatusTable(v)
	case *AppStoreVersionSubmissionCancelResult:
		return printAppStoreVersionSubmissionCancelTable(v)
	case *AppStoreVersionDetailResult:
		return printAppStoreVersionDetailTable(v)
	case *AppStoreVersionAttachBuildResult:
		return printAppStoreVersionAttachBuildTable(v)
	case *BuildBetaGroupsUpdateResult:
		return printBuildBetaGroupsUpdateTable(v)
	case *BetaTesterDeleteResult:
		return printBetaTesterDeleteResultTable(v)
	case *BetaTesterGroupsUpdateResult:
		return printBetaTesterGroupsUpdateResultTable(v)
	case *AppStoreVersionLocalizationDeleteResult:
		return printAppStoreVersionLocalizationDeleteResultTable(v)
	case *BetaTesterInvitationResult:
		return printBetaTesterInvitationResultTable(v)
	case *SandboxTesterDeleteResult:
		return printSandboxTesterDeleteResultTable(v)
	case *SandboxTesterClearHistoryResult:
		return printSandboxTesterClearHistoryResultTable(v)
	case *XcodeCloudRunResult:
		return printXcodeCloudRunResultTable(v)
	case *XcodeCloudStatusResult:
		return printXcodeCloudStatusResultTable(v)
	case *CiProductsResponse:
		return printCiProductsTable(v)
	case *CiWorkflowsResponse:
		return printCiWorkflowsTable(v)
	case *CiBuildRunsResponse:
		return printCiBuildRunsTable(v)
	case *CustomerReviewResponseResponse:
		return printCustomerReviewResponseTable(v)
	case *CustomerReviewResponseDeleteResult:
		return printCustomerReviewResponseDeleteResultTable(v)
	default:
		return PrintJSON(data)
	}
}

func compactWhitespace(input string) string {
	clean := sanitizeTerminal(input)
	return strings.Join(strings.Fields(clean), " ")
}

func escapeMarkdown(input string) string {
	clean := compactWhitespace(input)
	return strings.ReplaceAll(clean, "|", "\\|")
}

func feedbackHasScreenshots(resp *FeedbackResponse) bool {
	for _, item := range resp.Data {
		if len(item.Attributes.Screenshots) > 0 {
			return true
		}
	}
	return false
}

func formatScreenshotURLs(images []FeedbackScreenshotImage) string {
	if len(images) == 0 {
		return ""
	}
	urls := make([]string, 0, len(images))
	for _, image := range images {
		if strings.TrimSpace(image.URL) == "" {
			continue
		}
		urls = append(urls, image.URL)
	}
	return strings.Join(urls, ", ")
}

func printFeedbackTable(resp *FeedbackResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	hasScreenshots := feedbackHasScreenshots(resp)
	if hasScreenshots {
		fmt.Fprintln(w, "Created\tEmail\tComment\tScreenshots")
	} else {
		fmt.Fprintln(w, "Created\tEmail\tComment")
	}
	for _, item := range resp.Data {
		if hasScreenshots {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				sanitizeTerminal(item.Attributes.CreatedDate),
				sanitizeTerminal(item.Attributes.Email),
				compactWhitespace(item.Attributes.Comment),
				sanitizeTerminal(formatScreenshotURLs(item.Attributes.Screenshots)),
			)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			sanitizeTerminal(item.Attributes.CreatedDate),
			sanitizeTerminal(item.Attributes.Email),
			compactWhitespace(item.Attributes.Comment),
		)
	}
	return w.Flush()
}

func printCrashesTable(resp *CrashesResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Created\tEmail\tDevice\tOS\tComment")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			sanitizeTerminal(item.Attributes.CreatedDate),
			sanitizeTerminal(item.Attributes.Email),
			sanitizeTerminal(item.Attributes.DeviceModel),
			sanitizeTerminal(item.Attributes.OSVersion),
			compactWhitespace(item.Attributes.Comment),
		)
	}
	return w.Flush()
}

func printReviewsTable(resp *ReviewsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Created\tRating\tTerritory\tTitle")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
			sanitizeTerminal(item.Attributes.CreatedDate),
			item.Attributes.Rating,
			sanitizeTerminal(item.Attributes.Territory),
			compactWhitespace(item.Attributes.Title),
		)
	}
	return w.Flush()
}

func printFeedbackMarkdown(resp *FeedbackResponse) error {
	hasScreenshots := feedbackHasScreenshots(resp)
	if hasScreenshots {
		fmt.Fprintln(os.Stdout, "| Created | Email | Comment | Screenshots |")
		fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	} else {
		fmt.Fprintln(os.Stdout, "| Created | Email | Comment |")
		fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	}
	for _, item := range resp.Data {
		if hasScreenshots {
			fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s |\n",
				escapeMarkdown(item.Attributes.CreatedDate),
				escapeMarkdown(item.Attributes.Email),
				escapeMarkdown(item.Attributes.Comment),
				escapeMarkdown(formatScreenshotURLs(item.Attributes.Screenshots)),
			)
			continue
		}
		fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
			escapeMarkdown(item.Attributes.CreatedDate),
			escapeMarkdown(item.Attributes.Email),
			escapeMarkdown(item.Attributes.Comment),
		)
	}
	return nil
}

func printCrashesMarkdown(resp *CrashesResponse) error {
	fmt.Fprintln(os.Stdout, "| Created | Email | Device | OS | Comment |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s |\n",
			escapeMarkdown(item.Attributes.CreatedDate),
			escapeMarkdown(item.Attributes.Email),
			escapeMarkdown(item.Attributes.DeviceModel),
			escapeMarkdown(item.Attributes.OSVersion),
			escapeMarkdown(item.Attributes.Comment),
		)
	}
	return nil
}

func printReviewsMarkdown(resp *ReviewsResponse) error {
	fmt.Fprintln(os.Stdout, "| Created | Rating | Territory | Title |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %d | %s | %s |\n",
			escapeMarkdown(item.Attributes.CreatedDate),
			item.Attributes.Rating,
			escapeMarkdown(item.Attributes.Territory),
			escapeMarkdown(item.Attributes.Title),
		)
	}
	return nil
}

func printAppsTable(resp *AppsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tBundle ID\tSKU")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.Name),
			item.Attributes.BundleID,
			item.Attributes.SKU,
		)
	}
	return w.Flush()
}

func printAppStoreVersionLocalizationsTable(resp *AppStoreVersionLocalizationsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Locale\tWhats New\tKeywords")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			item.Attributes.Locale,
			compactWhitespace(item.Attributes.WhatsNew),
			compactWhitespace(item.Attributes.Keywords),
		)
	}
	return w.Flush()
}

func printAppInfoLocalizationsTable(resp *AppInfoLocalizationsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Locale\tName\tSubtitle\tPrivacy Policy URL")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			item.Attributes.Locale,
			compactWhitespace(item.Attributes.Name),
			compactWhitespace(item.Attributes.Subtitle),
			item.Attributes.PrivacyPolicyURL,
		)
	}
	return w.Flush()
}

func printBetaGroupsTable(resp *BetaGroupsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tInternal\tPublic Link Enabled\tPublic Link")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%t\t%t\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.Name),
			item.Attributes.IsInternalGroup,
			item.Attributes.PublicLinkEnabled,
			item.Attributes.PublicLink,
		)
	}
	return w.Flush()
}

func formatBetaTesterName(attr BetaTesterAttributes) string {
	first := strings.TrimSpace(attr.FirstName)
	last := strings.TrimSpace(attr.LastName)
	switch {
	case first == "" && last == "":
		return ""
	case first == "":
		return last
	case last == "":
		return first
	default:
		return first + " " + last
	}
}

func printBetaTestersTable(resp *BetaTestersResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tEmail\tName\tState\tInvite")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			item.Attributes.Email,
			compactWhitespace(formatBetaTesterName(item.Attributes)),
			string(item.Attributes.State),
			string(item.Attributes.InviteType),
		)
	}
	return w.Flush()
}

func printBetaTesterTable(resp *BetaTesterResponse) error {
	return printBetaTestersTable(&BetaTestersResponse{
		Data: []Resource[BetaTesterAttributes]{resp.Data},
	})
}

func printBuildsTable(resp *BuildsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Version\tUploaded\tProcessing\tExpired")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%t\n",
			item.Attributes.Version,
			item.Attributes.UploadedDate,
			item.Attributes.ProcessingState,
			item.Attributes.Expired,
		)
	}
	return w.Flush()
}

func printAppStoreVersionsTable(resp *AppStoreVersionsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tVersion\tPlatform\tState\tCreated")
	for _, item := range resp.Data {
		state := item.Attributes.AppVersionState
		if state == "" {
			state = item.Attributes.AppStoreState
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			item.Attributes.VersionString,
			string(item.Attributes.Platform),
			state,
			item.Attributes.CreatedDate,
		)
	}
	return w.Flush()
}

func printPreReleaseVersionsTable(resp *PreReleaseVersionsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tVersion\tPlatform")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.Version),
			string(item.Attributes.Platform),
		)
	}
	return w.Flush()
}

func printAppsMarkdown(resp *AppsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Name | Bundle ID | SKU |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s |\n",
			item.ID,
			escapeMarkdown(item.Attributes.Name),
			escapeMarkdown(item.Attributes.BundleID),
			escapeMarkdown(item.Attributes.SKU),
		)
	}
	return nil
}

func printAppStoreVersionsMarkdown(resp *AppStoreVersionsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Version | Platform | State | Created |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		state := item.Attributes.AppVersionState
		if state == "" {
			state = item.Attributes.AppStoreState
		}
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.VersionString),
			escapeMarkdown(string(item.Attributes.Platform)),
			escapeMarkdown(state),
			escapeMarkdown(item.Attributes.CreatedDate),
		)
	}
	return nil
}

func printPreReleaseVersionsMarkdown(resp *PreReleaseVersionsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Version | Platform |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.Version),
			escapeMarkdown(string(item.Attributes.Platform)),
		)
	}
	return nil
}

func printAppStoreVersionLocalizationsMarkdown(resp *AppStoreVersionLocalizationsResponse) error {
	fmt.Fprintln(os.Stdout, "| Locale | Whats New | Keywords |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
			escapeMarkdown(item.Attributes.Locale),
			escapeMarkdown(item.Attributes.WhatsNew),
			escapeMarkdown(item.Attributes.Keywords),
		)
	}
	return nil
}

func printAppInfoLocalizationsMarkdown(resp *AppInfoLocalizationsResponse) error {
	fmt.Fprintln(os.Stdout, "| Locale | Name | Subtitle | Privacy Policy URL |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s |\n",
			escapeMarkdown(item.Attributes.Locale),
			escapeMarkdown(item.Attributes.Name),
			escapeMarkdown(item.Attributes.Subtitle),
			escapeMarkdown(item.Attributes.PrivacyPolicyURL),
		)
	}
	return nil
}

func printBetaGroupsMarkdown(resp *BetaGroupsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Name | Internal | Public Link Enabled | Public Link |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %t | %t | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.Name),
			item.Attributes.IsInternalGroup,
			item.Attributes.PublicLinkEnabled,
			escapeMarkdown(item.Attributes.PublicLink),
		)
	}
	return nil
}

func printBetaTestersMarkdown(resp *BetaTestersResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Email | Name | State | Invite |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.Email),
			escapeMarkdown(formatBetaTesterName(item.Attributes)),
			escapeMarkdown(string(item.Attributes.State)),
			escapeMarkdown(string(item.Attributes.InviteType)),
		)
	}
	return nil
}

func printBetaTesterMarkdown(resp *BetaTesterResponse) error {
	return printBetaTestersMarkdown(&BetaTestersResponse{
		Data: []Resource[BetaTesterAttributes]{resp.Data},
	})
}

func printBuildsMarkdown(resp *BuildsResponse) error {
	fmt.Fprintln(os.Stdout, "| Version | Uploaded | Processing | Expired |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %t |\n",
			escapeMarkdown(item.Attributes.Version),
			escapeMarkdown(item.Attributes.UploadedDate),
			escapeMarkdown(item.Attributes.ProcessingState),
			item.Attributes.Expired,
		)
	}
	return nil
}

func printBuildUploadResultTable(result *BuildUploadResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Upload ID\tFile ID\tFile Name\tFile Size")
	fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
		result.UploadID,
		result.FileID,
		result.FileName,
		result.FileSize,
	)
	if err := w.Flush(); err != nil {
		return err
	}
	if len(result.Operations) == 0 {
		return nil
	}
	fmt.Fprintln(os.Stdout, "\nUpload Operations")
	opsWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(opsWriter, "Method\tURL\tLength\tOffset")
	for _, op := range result.Operations {
		fmt.Fprintf(opsWriter, "%s\t%s\t%d\t%d\n",
			op.Method,
			op.URL,
			op.Length,
			op.Offset,
		)
	}
	return opsWriter.Flush()
}

func printAppStoreVersionSubmissionTable(result *AppStoreVersionSubmissionResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Submission ID\tCreated Date")
	createdDate := ""
	if result.CreatedDate != nil {
		createdDate = *result.CreatedDate
	}
	fmt.Fprintf(w, "%s\t%s\n", result.SubmissionID, createdDate)
	return w.Flush()
}

func printAppStoreVersionSubmissionCreateTable(result *AppStoreVersionSubmissionCreateResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Submission ID\tVersion ID\tBuild ID\tCreated Date")
	createdDate := ""
	if result.CreatedDate != nil {
		createdDate = *result.CreatedDate
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
		result.SubmissionID,
		result.VersionID,
		result.BuildID,
		createdDate,
	)
	return w.Flush()
}

func printAppStoreVersionSubmissionStatusTable(result *AppStoreVersionSubmissionStatusResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Submission ID\tVersion ID\tVersion\tPlatform\tState\tCreated Date")
	createdDate := ""
	if result.CreatedDate != nil {
		createdDate = *result.CreatedDate
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
		result.ID,
		result.VersionID,
		result.VersionString,
		result.Platform,
		result.State,
		createdDate,
	)
	return w.Flush()
}

func printAppStoreVersionSubmissionCancelTable(result *AppStoreVersionSubmissionCancelResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Submission ID\tCancelled")
	fmt.Fprintf(w, "%s\t%t\n", result.ID, result.Cancelled)
	return w.Flush()
}

func printAppStoreVersionDetailTable(result *AppStoreVersionDetailResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Version ID\tVersion\tPlatform\tState\tBuild ID\tBuild Version\tSubmission ID")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		result.ID,
		result.VersionString,
		result.Platform,
		result.State,
		result.BuildID,
		result.BuildVersion,
		result.SubmissionID,
	)
	return w.Flush()
}

func printAppStoreVersionAttachBuildTable(result *AppStoreVersionAttachBuildResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Version ID\tBuild ID\tAttached")
	fmt.Fprintf(w, "%s\t%s\t%t\n", result.VersionID, result.BuildID, result.Attached)
	return w.Flush()
}

func printBuildBetaGroupsUpdateTable(result *BuildBetaGroupsUpdateResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Build ID\tGroup IDs\tAction")
	fmt.Fprintf(w, "%s\t%s\t%s\n",
		result.BuildID,
		strings.Join(result.GroupIDs, ", "),
		result.Action,
	)
	return w.Flush()
}

func printBuildUploadResultMarkdown(result *BuildUploadResult) error {
	fmt.Fprintln(os.Stdout, "| Upload ID | File ID | File Name | File Size |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %s | %d |\n",
		escapeMarkdown(result.UploadID),
		escapeMarkdown(result.FileID),
		escapeMarkdown(result.FileName),
		result.FileSize,
	)
	if len(result.Operations) == 0 {
		return nil
	}
	fmt.Fprintln(os.Stdout, "\n| Method | URL | Length | Offset |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, op := range result.Operations {
		fmt.Fprintf(os.Stdout, "| %s | %s | %d | %d |\n",
			escapeMarkdown(op.Method),
			escapeMarkdown(op.URL),
			op.Length,
			op.Offset,
		)
	}
	return nil
}

func printAppStoreVersionSubmissionMarkdown(result *AppStoreVersionSubmissionResult) error {
	fmt.Fprintln(os.Stdout, "| Submission ID | Created Date |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	createdDate := ""
	if result.CreatedDate != nil {
		createdDate = *result.CreatedDate
	}
	fmt.Fprintf(os.Stdout, "| %s | %s |\n",
		escapeMarkdown(result.SubmissionID),
		escapeMarkdown(createdDate),
	)
	return nil
}

func printAppStoreVersionSubmissionCreateMarkdown(result *AppStoreVersionSubmissionCreateResult) error {
	fmt.Fprintln(os.Stdout, "| Submission ID | Version ID | Build ID | Created Date |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	createdDate := ""
	if result.CreatedDate != nil {
		createdDate = *result.CreatedDate
	}
	fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s |\n",
		escapeMarkdown(result.SubmissionID),
		escapeMarkdown(result.VersionID),
		escapeMarkdown(result.BuildID),
		escapeMarkdown(createdDate),
	)
	return nil
}

func printAppStoreVersionSubmissionStatusMarkdown(result *AppStoreVersionSubmissionStatusResult) error {
	fmt.Fprintln(os.Stdout, "| Submission ID | Version ID | Version | Platform | State | Created Date |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- | --- |")
	createdDate := ""
	if result.CreatedDate != nil {
		createdDate = *result.CreatedDate
	}
	fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s | %s |\n",
		escapeMarkdown(result.ID),
		escapeMarkdown(result.VersionID),
		escapeMarkdown(result.VersionString),
		escapeMarkdown(result.Platform),
		escapeMarkdown(result.State),
		escapeMarkdown(createdDate),
	)
	return nil
}

func printAppStoreVersionSubmissionCancelMarkdown(result *AppStoreVersionSubmissionCancelResult) error {
	fmt.Fprintln(os.Stdout, "| Submission ID | Cancelled |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %t |\n",
		escapeMarkdown(result.ID),
		result.Cancelled,
	)
	return nil
}

func printAppStoreVersionDetailMarkdown(result *AppStoreVersionDetailResult) error {
	fmt.Fprintln(os.Stdout, "| Version ID | Version | Platform | State | Build ID | Build Version | Submission ID |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s | %s | %s |\n",
		escapeMarkdown(result.ID),
		escapeMarkdown(result.VersionString),
		escapeMarkdown(result.Platform),
		escapeMarkdown(result.State),
		escapeMarkdown(result.BuildID),
		escapeMarkdown(result.BuildVersion),
		escapeMarkdown(result.SubmissionID),
	)
	return nil
}

func printAppStoreVersionAttachBuildMarkdown(result *AppStoreVersionAttachBuildResult) error {
	fmt.Fprintln(os.Stdout, "| Version ID | Build ID | Attached |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %t |\n",
		escapeMarkdown(result.VersionID),
		escapeMarkdown(result.BuildID),
		result.Attached,
	)
	return nil
}

func printBuildBetaGroupsUpdateMarkdown(result *BuildBetaGroupsUpdateResult) error {
	fmt.Fprintln(os.Stdout, "| Build ID | Group IDs | Action |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
		escapeMarkdown(result.BuildID),
		escapeMarkdown(strings.Join(result.GroupIDs, ", ")),
		escapeMarkdown(result.Action),
	)
	return nil
}

func printLocalizationDownloadResultTable(result *LocalizationDownloadResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Locale\tPath")
	for _, file := range result.Files {
		fmt.Fprintf(w, "%s\t%s\n", file.Locale, file.Path)
	}
	return w.Flush()
}

func printLocalizationDownloadResultMarkdown(result *LocalizationDownloadResult) error {
	fmt.Fprintln(os.Stdout, "| Locale | Path |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	for _, file := range result.Files {
		fmt.Fprintf(os.Stdout, "| %s | %s |\n",
			escapeMarkdown(file.Locale),
			escapeMarkdown(file.Path),
		)
	}
	return nil
}

func printLocalizationUploadResultTable(result *LocalizationUploadResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Locale\tAction\tLocalization ID")
	for _, item := range result.Results {
		fmt.Fprintf(w, "%s\t%s\t%s\n", item.Locale, item.Action, item.LocalizationID)
	}
	return w.Flush()
}

func printLocalizationUploadResultMarkdown(result *LocalizationUploadResult) error {
	fmt.Fprintln(os.Stdout, "| Locale | Action | Localization ID |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	for _, item := range result.Results {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
			escapeMarkdown(item.Locale),
			escapeMarkdown(item.Action),
			escapeMarkdown(item.LocalizationID),
		)
	}
	return nil
}

func printAppStoreVersionLocalizationDeleteResultTable(result *AppStoreVersionLocalizationDeleteResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDeleted")
	fmt.Fprintf(w, "%s\t%t\n",
		result.ID,
		result.Deleted,
	)
	return w.Flush()
}

func printAppStoreVersionLocalizationDeleteResultMarkdown(result *AppStoreVersionLocalizationDeleteResult) error {
	fmt.Fprintln(os.Stdout, "| ID | Deleted |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %t |\n",
		escapeMarkdown(result.ID),
		result.Deleted,
	)
	return nil
}

func printBetaTesterDeleteResultTable(result *BetaTesterDeleteResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tEmail\tDeleted")
	fmt.Fprintf(w, "%s\t%s\t%t\n",
		result.ID,
		result.Email,
		result.Deleted,
	)
	return w.Flush()
}

func printBetaTesterDeleteResultMarkdown(result *BetaTesterDeleteResult) error {
	fmt.Fprintln(os.Stdout, "| ID | Email | Deleted |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %t |\n",
		escapeMarkdown(result.ID),
		escapeMarkdown(result.Email),
		result.Deleted,
	)
	return nil
}

func printBetaTesterGroupsUpdateResultTable(result *BetaTesterGroupsUpdateResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Tester ID\tGroup IDs\tAction")
	fmt.Fprintf(w, "%s\t%s\t%s\n",
		result.TesterID,
		strings.Join(result.GroupIDs, ","),
		result.Action,
	)
	return w.Flush()
}

func printBetaTesterGroupsUpdateResultMarkdown(result *BetaTesterGroupsUpdateResult) error {
	fmt.Fprintln(os.Stdout, "| Tester ID | Group IDs | Action |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
		escapeMarkdown(result.TesterID),
		escapeMarkdown(strings.Join(result.GroupIDs, ",")),
		escapeMarkdown(result.Action),
	)
	return nil
}

func printBetaTesterInvitationResultTable(result *BetaTesterInvitationResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Invitation ID\tTester ID\tApp ID\tEmail")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
		result.InvitationID,
		result.TesterID,
		result.AppID,
		result.Email,
	)
	return w.Flush()
}

func printBetaTesterInvitationResultMarkdown(result *BetaTesterInvitationResult) error {
	fmt.Fprintln(os.Stdout, "| Invitation ID | Tester ID | App ID | Email |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s |\n",
		escapeMarkdown(result.InvitationID),
		escapeMarkdown(result.TesterID),
		escapeMarkdown(result.AppID),
		escapeMarkdown(result.Email),
	)
	return nil
}
