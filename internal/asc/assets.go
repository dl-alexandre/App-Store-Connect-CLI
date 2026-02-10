package asc

import "slices"

// AppScreenshotSetAttributes describes a screenshot set resource.
type AppScreenshotSetAttributes struct {
	ScreenshotDisplayType string `json:"screenshotDisplayType"`
}

// AppScreenshotAttributes describes a screenshot asset resource.
type AppScreenshotAttributes struct {
	FileSize           int64               `json:"fileSize"`
	FileName           string              `json:"fileName"`
	SourceFileChecksum string              `json:"sourceFileChecksum,omitempty"`
	ImageAsset         *ImageAsset         `json:"imageAsset,omitempty"`
	AssetToken         string              `json:"assetToken,omitempty"`
	AssetType          string              `json:"assetType,omitempty"`
	UploadOperations   []UploadOperation   `json:"uploadOperations,omitempty"`
	AssetDeliveryState *AssetDeliveryState `json:"assetDeliveryState,omitempty"`
}

// ImageAsset describes an image asset.
type ImageAsset struct {
	TemplateURL string `json:"templateUrl"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

// AssetDeliveryState describes the delivery state of an asset.
type AssetDeliveryState struct {
	State  string        `json:"state"`
	Errors []ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail describes an asset error detail.
type ErrorDetail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// AppPreviewSetAttributes describes a preview set resource.
type AppPreviewSetAttributes struct {
	PreviewType string `json:"previewType"`
}

// AppPreviewAttributes describes a preview asset resource.
type AppPreviewAttributes struct {
	FileSize             int64               `json:"fileSize"`
	FileName             string              `json:"fileName"`
	SourceFileChecksum   string              `json:"sourceFileChecksum,omitempty"`
	PreviewFrameTimeCode string              `json:"previewFrameTimeCode,omitempty"`
	MimeType             string              `json:"mimeType,omitempty"`
	VideoURL             string              `json:"videoUrl,omitempty"`
	PreviewImage         *ImageAsset         `json:"previewImage,omitempty"`
	UploadOperations     []UploadOperation   `json:"uploadOperations,omitempty"`
	AssetDeliveryState   *AssetDeliveryState `json:"assetDeliveryState,omitempty"`
}

// Response types
type (
	AppScreenshotSetsResponse = Response[AppScreenshotSetAttributes]
	AppScreenshotSetResponse  = SingleResponse[AppScreenshotSetAttributes]
	AppScreenshotsResponse    = Response[AppScreenshotAttributes]
	AppScreenshotResponse     = SingleResponse[AppScreenshotAttributes]
	AppPreviewSetsResponse    = Response[AppPreviewSetAttributes]
	AppPreviewSetResponse     = SingleResponse[AppPreviewSetAttributes]
	AppPreviewsResponse       = Response[AppPreviewAttributes]
	AppPreviewResponse        = SingleResponse[AppPreviewAttributes]
)

// Valid screenshot display types for validation.
var ValidScreenshotDisplayTypes = []string{
	"APP_IPHONE_69",
	"APP_IPHONE_67",
	"APP_IPHONE_65",
	"APP_IPHONE_61",
	"APP_IPHONE_58",
	"APP_IPHONE_55",
	"APP_IPHONE_47",
	"APP_IPHONE_40",
	"APP_IPHONE_35",
	"APP_IPAD_PRO_3GEN_129",
	"APP_IPAD_PRO_3GEN_11",
	"APP_IPAD_PRO_129",
	"APP_IPAD_105",
	"APP_IPAD_97",
	"APP_DESKTOP",
	"APP_WATCH_ULTRA",
	"APP_WATCH_SERIES_10",
	"APP_WATCH_SERIES_7",
	"APP_WATCH_SERIES_4",
	"APP_WATCH_SERIES_3",
	"APP_APPLE_TV",
	"APP_APPLE_VISION_PRO",
	"IMESSAGE_APP_IPHONE_69",
	"IMESSAGE_APP_IPHONE_67",
	"IMESSAGE_APP_IPHONE_61",
	"IMESSAGE_APP_IPHONE_65",
	"IMESSAGE_APP_IPHONE_58",
	"IMESSAGE_APP_IPHONE_55",
	"IMESSAGE_APP_IPHONE_47",
	"IMESSAGE_APP_IPHONE_40",
	"IMESSAGE_APP_IPAD_PRO_3GEN_129",
	"IMESSAGE_APP_IPAD_PRO_3GEN_11",
	"IMESSAGE_APP_IPAD_PRO_129",
	"IMESSAGE_APP_IPAD_105",
	"IMESSAGE_APP_IPAD_97",
}

// Valid preview types for validation.
var ValidPreviewTypes = []string{
	"IPHONE_67",
	"IPHONE_65",
	"IPHONE_61",
	"IPHONE_58",
	"IPHONE_55",
	"IPHONE_47",
	"IPHONE_40",
	"IPHONE_35",
	"IPAD_PRO_3GEN_129",
	"IPAD_PRO_3GEN_11",
	"IPAD_PRO_129",
	"IPAD_105",
	"IPAD_97",
	"DESKTOP",
	"APPLE_TV",
	"APPLE_VISION_PRO",
}

// IsValidScreenshotDisplayType checks if a screenshot display type is supported.
func IsValidScreenshotDisplayType(value string) bool {
	return slices.Contains(ValidScreenshotDisplayTypes, value)
}

// IsValidPreviewType checks if a preview type is supported.
func IsValidPreviewType(value string) bool {
	return slices.Contains(ValidPreviewTypes, value)
}
