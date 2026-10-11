package zephior

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testClient(t *testing.T, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			data, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(data))
			var body map[string]any
			_ = json.Unmarshal(data, &body)
			if body["method"] == "initialize" {
				_ = json.NewEncoder(w).Encode(map[string]any{"result": map[string]any{"protocolVersion": "2025-06-18"}})
				return
			}
			if body["method"] == "notifications/initialized" {
				w.WriteHeader(202)
				return
			}
		}
		handler(w, r)
	}))
	return &Client{BaseURL: server.URL, HTTPClient: server.Client()}, server.Close
}

func TestGetSiteOverview(t *testing.T) {
	client, closeServer := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/site" || r.URL.Query().Get("locale") != "de" {
			t.Fatalf("unexpected request %s", r.URL.String())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"name": "Zephior", "locale": "de"})
	})
	defer closeServer()
	result, err := client.GetSiteOverview(context.Background(), "de")
	if err != nil {
		t.Fatal(err)
	}
	if result["locale"] != "de" {
		t.Fatalf("unexpected locale %v", result["locale"])
	}
}

func TestPublicTenderFilters(t *testing.T) {
	client, closeServer := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tenders" || r.URL.Query().Get("countries") != "CH,DE" || !r.URL.Query().Has("statuses") || r.URL.Query().Get("statuses") != "" || r.URL.Query().Get("q") != "bridge" {
			t.Errorf("unexpected combined query %s", r.URL)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"items": []any{}, "total": 0, "page": 2, "limit": 20})
	})
	defer closeServer()
	result, err := client.SearchPublicTenders(context.Background(), TenderSearch{Countries: []string{"CH", "DE"}, Statuses: []string{}, Query: "bridge", Page: 2})
	if err != nil || result.Page != 2 || result.Limit != 20 {
		t.Fatalf("unexpected result %v: %v", result, err)
	}
}

func TestCreateZivaPricingEstimate(t *testing.T) {
	client, closeServer := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mcp/product" || r.Method != http.MethodPost {
			t.Errorf("unexpected pricing request %s %s", r.Method, r.URL)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		params := request["params"].(map[string]any)
		if params["name"] != "create_ziva_pricing_estimate" {
			t.Fatalf("unexpected MCP tool %v", params["name"])
		}
		arguments := params["arguments"].(map[string]any)
		if arguments["projectKeys"] != float64(1) {
			t.Fatalf("unexpected project keys %v", arguments["projectKeys"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]any{
				"structuredContent": map[string]any{"amount": 1149, "currency": "CHF"},
			},
		})
	})
	defer closeServer()
	result, err := client.CreateZivaPricingEstimate(context.Background(), PricingEstimateInput{ProjectKeys: 1, Currency: "CHF", ProjectsBillingPeriod: "monthly"})
	if err != nil {
		t.Fatal(err)
	}
	if result["amount"] != float64(1149) {
		t.Fatalf("unexpected amount %v", result["amount"])
	}
}

func TestCallMCPToolAtSelectsPublicServer(t *testing.T) {
	for _, path := range []string{"/mcp", "/mcp/product"} {
		t.Run(path, func(t *testing.T) {
			client, closeServer := testClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != path || r.Method != http.MethodPost {
					t.Errorf("unexpected MCP request %s %s", r.Method, r.URL)
				}
				var request struct {
					JSONRPC string `json:"jsonrpc"`
					Method  string `json:"method"`
					Params  struct {
						Name      string         `json:"name"`
						Arguments map[string]any `json:"arguments"`
					} `json:"params"`
				}
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Error(err)
					http.Error(w, "invalid request", http.StatusBadRequest)
					return
				}
				if request.JSONRPC != "2.0" || request.Method != "tools/call" || request.Params.Name != "test_tool" || request.Params.Arguments["locale"] != "fr" {
					t.Errorf("unexpected MCP payload %+v", request)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"jsonrpc": "2.0", "id": 1,
					"result": map[string]any{"structuredContent": map[string]any{"server": path}},
				})
			})
			defer closeServer()
			result, err := client.CallMCPToolAt(context.Background(), path, "test_tool", map[string]any{"locale": "fr"})
			if err != nil || result["server"] != path {
				t.Fatalf("unexpected result %v, error %v", result, err)
			}
		})
	}
}

func TestLegacyMCPAndTypedPricingStayOnProductServer(t *testing.T) {
	client, closeServer := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mcp/product" {
			t.Errorf("unexpected legacy or pricing endpoint %s", r.URL)
		}
		var request struct {
			Params struct {
				Arguments map[string]any `json:"arguments"`
			} `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		if request.Params.Arguments == nil {
			t.Error("MCP arguments must be an object even for nil input")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{"structuredContent": map[string]any{"ok": true}},
		})
	})
	defer closeServer()
	if _, err := client.CallMCPTool(context.Background(), "legacy_tool", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetZivaPricingDocument(context.Background(), ""); err != nil {
		t.Fatal(err)
	}
}

func TestCallMCPToolAtRejectsInvalidTargetsBeforeRequest(t *testing.T) {
	// No HTTP client: invalid input must return before attempting a request.
	client := &Client{BaseURL: DefaultOrigin}
	for _, path := range []string{"", "/api/v1/search", "/mcp/", "/mcp/../private", "/mcp?x=1", "https://example.com/mcp", "//example.com/mcp"} {
		if _, err := client.CallMCPToolAt(context.Background(), path, "test_tool", nil); err == nil {
			t.Errorf("accepted invalid target %q", path)
		}
	}
	for _, name := range []string{"", " \t\n"} {
		if _, err := client.CallMCPToolAt(context.Background(), "/mcp", name, nil); err == nil {
			t.Errorf("accepted empty tool name %q", name)
		}
	}
}

func TestCallMCPToolAtPreservesErrors(t *testing.T) {
	for _, path := range []string{"/mcp", "/mcp/product"} {
		t.Run(path, func(t *testing.T) {
			client, closeServer := testClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != path {
					t.Errorf("unexpected MCP path %s", r.URL.Path)
				}
				if r.URL.Query().Get("unused") != "" {
					t.Error("unexpected query")
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]any{"code": -32601, "message": "Unknown tool"},
				})
			})
			defer closeServer()
			if _, err := client.CallMCPToolAt(context.Background(), path, "unknown", nil); err == nil || !strings.Contains(err.Error(), "-32601: Unknown tool") {
				t.Fatalf("lost MCP error: %v", err)
			}
		})
	}
	client, closeServer := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": "rate_limit_exceeded", "title": "Rate limit exceeded", "resolution": "Retry later",
		})
	})
	defer closeServer()
	_, err := client.CallMCPToolAt(context.Background(), "/mcp", "test_tool", nil)
	var publicErr *PublicAPIError
	if !errors.As(err, &publicErr) || publicErr.Status != http.StatusTooManyRequests || publicErr.Resolution != "Retry later" {
		t.Fatalf("lost HTTP problem: %v", err)
	}
}
