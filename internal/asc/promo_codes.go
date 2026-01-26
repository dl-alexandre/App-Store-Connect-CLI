package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// PromoCodeProductType represents the product type for promo codes.
type PromoCodeProductType string

const (
	PromoCodeProductTypeApp          PromoCodeProductType = "APP"
	PromoCodeProductTypeSubscription PromoCodeProductType = "SUBSCRIPTION"
)

// PromoCodeAttributes describes a promo code resource.
type PromoCodeAttributes struct {
	Code           string               `json:"code,omitempty"`
	ExpiresDate    string               `json:"expiresDate,omitempty"`
	ExpirationDate string               `json:"expirationDate,omitempty"`
	IsExpired      bool                 `json:"isExpired,omitempty"`
	IsUsed         bool                 `json:"isUsed,omitempty"`
	ProductType    PromoCodeProductType `json:"productType,omitempty"`
	Type           PromoCodeProductType `json:"type,omitempty"`
}

// PromoCodesResponse is the response from promo codes list/create endpoints.
type PromoCodesResponse = Response[PromoCodeAttributes]

// PromoCodeResponse is the response from promo code detail endpoint.
type PromoCodeResponse = SingleResponse[PromoCodeAttributes]

// PromoCodeFundamentalMetadata describes optional prefix configuration.
type PromoCodeFundamentalMetadata struct {
	CodePrefix string `json:"codePrefix,omitempty"`
}

// PromoCodeCreateAttributes describes attributes for generating promo codes.
type PromoCodeCreateAttributes struct {
	ProductType         PromoCodeProductType          `json:"productType"`
	Quantity            int                           `json:"quantity"`
	FundamentalMetadata *PromoCodeFundamentalMetadata `json:"fundamentalMetadata,omitempty"`
}

// PromoCodeCreateRelationships describes relationships for promo code creation.
type PromoCodeCreateRelationships struct {
	App Relationship `json:"app"`
}

// PromoCodeCreateData is the data portion of a promo code create request.
type PromoCodeCreateData struct {
	Type          ResourceType                 `json:"type"`
	Attributes    PromoCodeCreateAttributes    `json:"attributes"`
	Relationships PromoCodeCreateRelationships `json:"relationships"`
}

// PromoCodeCreateRequest is a request to generate promo codes.
type PromoCodeCreateRequest struct {
	Data PromoCodeCreateData `json:"data"`
}

// GetPromoCodes retrieves promo codes for an app.
func (c *Client) GetPromoCodes(ctx context.Context, appID string, opts ...PromoCodesOption) (*PromoCodesResponse, error) {
	appID = strings.TrimSpace(appID)
	query := &promoCodesQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/apps/%s/promoCodes", appID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("promo codes: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildPromoCodesQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response PromoCodesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetPromoCode retrieves a promo code by ID.
func (c *Client) GetPromoCode(ctx context.Context, promoCodeID string) (*PromoCodeResponse, error) {
	promoCodeID = strings.TrimSpace(promoCodeID)
	path := fmt.Sprintf("/v1/promoCodes/%s", promoCodeID)

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response PromoCodeResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreatePromoCodes generates new promo codes.
func (c *Client) CreatePromoCodes(ctx context.Context, req PromoCodeCreateRequest) (*PromoCodesResponse, error) {
	body, err := BuildRequestBody(req)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/promoCodes", body)
	if err != nil {
		return nil, err
	}

	var response PromoCodesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
