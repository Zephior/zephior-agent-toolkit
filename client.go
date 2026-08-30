package zephior

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const DefaultOrigin = "https://zephior.com"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type TenderSearchOptions struct {
	Countries           []string `json:"countries,omitempty"`
	Portals             []string `json:"portals,omitempty"`
	Statuses            []string `json:"statuses,omitempty"`
	ClassificationCodes []string `json:"classification_codes,omitempty"`
	PublishedAfter      string   `json:"published_after,omitempty"`
	DeadlineAfter       string   `json:"deadline_after,omitempty"`
	DeadlineBefore      string   `json:"deadline_before,omitempty"`
	Limit               int      `json:"limit,omitempty"`
}

type PublicAPIError struct {
	Status     int
	Code       any
	Title      string
	Detail     string
	Resolution string
}

func (e *PublicAPIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("zephior: %s: %s", e.Title, e.Detail)
	}
	return fmt.Sprintf("zephior: %s (HTTP %d)", e.Title, e.Status)
}

func NewClient() *Client {
	return &Client{BaseURL: DefaultOrigin, HTTPClient: http.DefaultClient}
}

func (c *Client) endpoint(path string) string {
	base := strings.TrimRight(c.BaseURL, "/")
	return base + "/" + strings.TrimLeft(path, "/")
}

func (c *Client) doJSON(ctx context.Context, method, endpoint string, body any, output any) error {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json, text/event-stream")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		problem := struct {
			Code       any    `json:"code"`
			Title      string `json:"title"`
			Detail     string `json:"detail"`
			Resolution string `json:"resolution"`
		}{}
		_ = json.Unmarshal(data, &problem)
		if problem.Title == "" {
			problem.Title = http.StatusText(response.StatusCode)
		}
		return &PublicAPIError{
			Status: response.StatusCode, Code: problem.Code, Title: problem.Title,
			Detail: problem.Detail, Resolution: problem.Resolution,
		}
	}
	if err := json.Unmarshal(data, output); err != nil {
		return fmt.Errorf("zephior: decode response: %w", err)
	}
	return nil
}

func (c *Client) GetSiteOverview(ctx context.Context, locale string) (map[string]any, error) {
	if locale == "" {
		locale = "en"
	}
	result := map[string]any{}
	err := c.doJSON(
		ctx,
		http.MethodGet,
		c.endpoint("/api/v1/site")+"?locale="+url.QueryEscape(locale),
		nil,
		&result,
	)
	return result, err
}

func (c *Client) SearchResources(
	ctx context.Context,
	query, locale string,
	limit int,
) (map[string]any, error) {
	if len(strings.TrimSpace(query)) < 2 {
		return nil, fmt.Errorf("zephior: query must contain at least two characters")
	}
	if locale == "" {
		locale = "en"
	}
	if limit == 0 {
		limit = 5
	}
	values := url.Values{
		"q":      {strings.TrimSpace(query)},
		"locale": {locale},
		"limit":  {fmt.Sprintf("%d", limit)},
	}
	result := map[string]any{}
	err := c.doJSON(ctx, http.MethodGet, c.endpoint("/api/v1/search")+"?"+values.Encode(), nil, &result)
	return result, err
}

func (c *Client) SearchLiveTenders(
	ctx context.Context,
	query string,
	options TenderSearchOptions,
) (map[string]any, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 || len(query) > 500 {
		return nil, fmt.Errorf("zephior: query must contain between two and 500 characters")
	}
	if options.Limit == 0 {
		options.Limit = 20
	}
	if options.Limit < 1 || options.Limit > 20 {
		return nil, fmt.Errorf("zephior: limit must be between 1 and 20")
	}
	if options.Statuses == nil {
		options.Statuses = []string{"open"}
	}
	payload := struct {
		Query string `json:"query"`
		TenderSearchOptions
	}{Query: query, TenderSearchOptions: options}
	result := map[string]any{}
	err := c.doJSON(ctx, http.MethodPost, c.endpoint("/api/tenders/search"), payload, &result)
	return result, err
}

func (c *Client) GetLatestLiveTenders(ctx context.Context, limit int) (map[string]any, error) {
	if limit == 0 {
		limit = 10
	}
	if limit < 1 || limit > 20 {
		return nil, fmt.Errorf("zephior: limit must be between 1 and 20")
	}
	result := map[string]any{}
	err := c.doJSON(
		ctx,
		http.MethodGet,
		c.endpoint("/api/tenders/latest")+"?limit="+url.QueryEscape(strconv.Itoa(limit)),
		nil,
		&result,
	)
	return result, err
}

func (c *Client) GetTenderUniverseStats(ctx context.Context) (map[string]any, error) {
	result := map[string]any{}
	err := c.doJSON(ctx, http.MethodGet, c.endpoint("/api/tenders/stats"), nil, &result)
	return result, err
}

func (c *Client) TenderUpdateStreamURL() string {
	return c.endpoint("/api/tenders/stream")
}

func (c *Client) CallMCPTool(
	ctx context.Context,
	name string,
	arguments map[string]any,
) (map[string]any, error) {
	request := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      name,
			"arguments": arguments,
		},
	}
	response := struct {
		Result struct {
			StructuredContent map[string]any `json:"structuredContent"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    any    `json:"data"`
		} `json:"error"`
	}{}
	if err := c.doJSON(ctx, http.MethodPost, c.endpoint("/mcp/product"), request, &response); err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("zephior MCP %d: %s", response.Error.Code, response.Error.Message)
	}
	return response.Result.StructuredContent, nil
}

func (c *Client) GetZivaPricingDocument(
	ctx context.Context,
	currency string,
) (map[string]any, error) {
	if currency == "" {
		currency = "CHF"
	}
	return c.CallMCPTool(ctx, "get_ziva_pricing_document", map[string]any{"currency": currency})
}

func (c *Client) CreateZivaPricingEstimate(
	ctx context.Context,
	credits int,
	currency, billingPeriod string,
) (map[string]any, error) {
	return c.CallMCPTool(ctx, "create_ziva_pricing_estimate", map[string]any{
		"credits": credits, "currency": currency, "billingPeriod": billingPeriod,
	})
}
