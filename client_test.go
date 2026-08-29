package zephior

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(t *testing.T, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
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

func TestCreateZivaPricingEstimate(t *testing.T) {
	client, closeServer := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		params := request["params"].(map[string]any)
		if params["name"] != "create_ziva_pricing_estimate" {
			t.Fatalf("unexpected MCP tool %v", params["name"])
		}
		arguments := params["arguments"].(map[string]any)
		if arguments["credits"] != float64(1149) {
			t.Fatalf("unexpected credits %v", arguments["credits"])
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
	result, err := client.CreateZivaPricingEstimate(context.Background(), 1149, "CHF", "monthly")
	if err != nil {
		t.Fatal(err)
	}
	if result["amount"] != float64(1149) {
		t.Fatalf("unexpected amount %v", result["amount"])
	}
}
