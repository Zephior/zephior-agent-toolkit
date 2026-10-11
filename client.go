package zephior

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
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
	return &Client{BaseURL: DefaultOrigin, HTTPClient: &http.Client{Timeout: 15 * time.Second}}
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
	if strings.HasPrefix(req.URL.Path, "/mcp") {
		req.Header.Set("MCP-Protocol-Version", "2025-06-18")
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
	data, err := io.ReadAll(io.LimitReader(response.Body, 8*1024*1024+1))
	if err != nil {
		return err
	}
	if len(data) > 8*1024*1024 {
		return fmt.Errorf("zephior: response exceeds 8 MiB")
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
	if strings.Contains(response.Header.Get("Content-Type"), "text/event-stream") {
		found := false
		for _, block := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n\n") {
			var lines []string
			for _, line := range strings.Split(block, "\n") {
				if strings.HasPrefix(line, "data:") {
					lines = append(lines, strings.TrimPrefix(line, "data:"))
				}
			}
			event := []byte(strings.Join(lines, "\n"))
			var envelope map[string]any
			if json.Unmarshal(event, &envelope) == nil && envelope["id"] != nil {
				data = event
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("zephior: missing JSON-RPC response in event stream")
		}
	}
	if output == nil && len(data) == 0 {
		return nil
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
	initRequest := map[string]any{"jsonrpc": "2.0", "id": 0, "method": "initialize", "params": map[string]any{"protocolVersion": "2025-06-18", "capabilities": map[string]any{}, "clientInfo": map[string]any{"name": "zephior-agent-toolkit", "version": "0.4.0"}}}
	var handshake struct {
		Result struct {
			ProtocolVersion string `json:"protocolVersion"`
		} `json:"result"`
	}
	if err := c.doJSON(ctx, http.MethodPost, c.endpoint(path), initRequest, &handshake); err != nil {
		return nil, err
	}
	if handshake.Result.ProtocolVersion != "2025-06-18" {
		return nil, fmt.Errorf("zephior: unsupported MCP protocol negotiated")
	}
	if err := c.doJSON(ctx, http.MethodPost, c.endpoint(path), map[string]any{"jsonrpc": "2.0", "method": "notifications/initialized"}, nil); err != nil {
		return nil, err
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
			IsError           bool           `json:"isError"`
			Content           []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
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
	if response.Result.IsError {
		return nil, fmt.Errorf("zephior: public tool could not complete: %v", response.Result.Content)
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

type PricingEstimateInput struct {
	Currency              string `json:"currency"`
	IntelligenceCoverage  string `json:"intelligenceCoverage,omitempty"`
	ProjectKeys           int    `json:"projectKeys,omitempty"`
	ProjectsBillingPeriod string `json:"projectsBillingPeriod,omitempty"`
}

func (c *Client) CreateZivaPricingEstimate(ctx context.Context, input PricingEstimateInput) (map[string]any, error) {
	if input.IntelligenceCoverage == "" && input.ProjectKeys == 0 {
		return nil, fmt.Errorf("zephior: select intelligence coverage, project keys, or both")
	}
	if input.Currency == "" {
		input.Currency = "CHF"
	}
	encoded, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var args map[string]any
	if err := json.Unmarshal(encoded, &args); err != nil {
		return nil, err
	}
	return c.CallMCPToolAt(ctx, "/mcp/product", "create_ziva_pricing_estimate", args)
}

type TenderSearch struct {
	Query      string   `json:"q,omitempty"`
	Countries  []string `json:"countries,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Statuses   []string `json:"statuses"`
	Buyer      string   `json:"buyer,omitempty"`
	Sort       string   `json:"sort,omitempty"`
	Page       int      `json:"page,omitempty"`
	Locale     string   `json:"locale,omitempty"`
}
type PublicTender struct {
	ID              string `json:"id"`
	URL             string `json:"url"`
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	Locale          string `json:"locale"`
	Status          string `json:"status"`
	RecordedStatus  string `json:"recordedStatus"`
	Country         string `json:"country"`
	ModifiedAt      string `json:"modifiedAt"`
	SourceCheckedAt string `json:"sourceCheckedAt"`
	Revision        int    `json:"revision"`
	Stale           bool   `json:"stale"`
	Deadline        struct {
		At        *string `json:"at"`
		Date      *string `json:"date"`
		Precision string  `json:"precision"`
	} `json:"deadline"`
	OriginalSource map[string]any   `json:"originalSource"`
	Sources        []map[string]any `json:"sources"`
	Content        map[string]any   `json:"content"`
	Facts          map[string]any   `json:"facts"`
	History        []map[string]any `json:"history"`
}
type TenderResults struct {
	Items    []PublicTender `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	Limit    int            `json:"limit"`
	NextPage *int           `json:"nextPage"`
	URL      string         `json:"url"`
	Locale   string         `json:"locale"`
	Filters  map[string]any `json:"filters"`
}

func (c *Client) SearchPublicTenders(ctx context.Context, input TenderSearch) (TenderResults, error) {
	query := url.Values{}
	for key, values := range map[string][]string{"countries": input.Countries, "categories": input.Categories, "statuses": input.Statuses} {
		if values != nil {
			query.Set(key, strings.Join(values, ","))
		}
	}
	for key, value := range map[string]string{"q": input.Query, "buyer": input.Buyer, "sort": input.Sort, "locale": input.Locale} {
		if value != "" {
			query.Set(key, value)
		}
	}
	if input.Page != 0 {
		query.Set("page", fmt.Sprint(input.Page))
	}
	var result TenderResults
	err := c.doJSON(ctx, http.MethodGet, c.endpoint("/api/v1/tenders")+"?"+query.Encode(), nil, &result)
	return result, err
}
func (c *Client) GetPublicTender(ctx context.Context, id, locale string) (PublicTender, error) {
	var result PublicTender
	if !regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(id) {
		return result, fmt.Errorf("zephior: a public notice UUID is required")
	}
	if locale == "" {
		locale = "en"
	}
	err := c.doJSON(ctx, http.MethodGet, c.endpoint("/api/v1/tenders/"+id)+"?locale="+url.QueryEscape(locale), nil, &result)
	return result, err
}
func (c *Client) GetPublicTenderFacets(ctx context.Context, locale string) (map[string]any, error) {
	if locale == "" {
		locale = "en"
	}
	result := map[string]any{}
	err := c.doJSON(ctx, http.MethodGet, c.endpoint("/api/v1/tenders/facets")+"?locale="+url.QueryEscape(locale), nil, &result)
	return result, err
}
