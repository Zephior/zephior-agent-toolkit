package zephior

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const DefaultOrigin = "https://zephior.com"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
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

// CallMCPTool calls the product-decision server for backward compatibility.
// Deprecated: use CallMCPToolAt with an explicit /mcp or /mcp/product path.
func (c *Client) CallMCPTool(
	ctx context.Context,
	name string,
	arguments map[string]any,
) (map[string]any, error) {
	return c.CallMCPToolAt(ctx, "/mcp/product", name, arguments)
}

// CallMCPToolAt calls one of the two public MCP servers on the client's BaseURL.
// Use /mcp for public-resource research, or /mcp/product
// for published pricing and service selection. Other paths and URLs are rejected.
func (c *Client) CallMCPToolAt(
	ctx context.Context,
	path, name string,
	arguments map[string]any,
) (map[string]any, error) {
	if path != "/mcp" && path != "/mcp/product" {
		return nil, fmt.Errorf("zephior: MCP path must be /mcp or /mcp/product")
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("zephior: MCP tool name is required")
	}
	if arguments == nil {
		arguments = map[string]any{}
	}
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
	if err := c.doJSON(ctx, http.MethodPost, c.endpoint(path), request, &response); err != nil {
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
	return c.CallMCPToolAt(ctx, "/mcp/product", "get_ziva_pricing_document", map[string]any{"currency": currency})
}

func (c *Client) CreateZivaPricingEstimate(
	ctx context.Context,
	credits int,
	currency, billingPeriod string,
) (map[string]any, error) {
	return c.CallMCPToolAt(ctx, "/mcp/product", "create_ziva_pricing_estimate", map[string]any{
		"credits": credits, "currency": currency, "billingPeriod": billingPeriod,
	})
}
