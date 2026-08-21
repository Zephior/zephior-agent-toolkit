package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	zephior "github.com/tonykzeph/zephior-agent-toolkit"
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:
  zephior site [locale]
  zephior search <query> [locale] [limit]
  zephior pricing [CHF|USD|EUR]
  zephior estimate <Starter|Business|Advanced|Enterprise> <CHF|USD|EUR> <monthly|annual>

The CLI calls Zephior's public, read-only interfaces. Never submit confidential material.`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	ctx := context.Background()
	client := zephior.NewClient()
	var (
		result map[string]any
		err    error
	)
	switch os.Args[1] {
	case "site":
		locale := "en"
		if len(os.Args) > 2 {
			locale = os.Args[2]
		}
		result, err = client.GetSiteOverview(ctx, locale)
	case "search":
		if len(os.Args) < 3 {
			usage()
			os.Exit(2)
		}
		locale, limit := "en", 5
		if len(os.Args) > 3 {
			locale = os.Args[3]
		}
		if len(os.Args) > 4 {
			limit, err = strconv.Atoi(os.Args[4])
			if err != nil {
				fmt.Fprintln(os.Stderr, "limit must be an integer")
				os.Exit(2)
			}
		}
		result, err = client.SearchResources(ctx, strings.TrimSpace(os.Args[2]), locale, limit)
	case "pricing":
		currency := "CHF"
		if len(os.Args) > 2 {
			currency = os.Args[2]
		}
		result, err = client.GetZivaPricingDocument(ctx, currency)
	case "estimate":
		if len(os.Args) != 5 {
			usage()
			os.Exit(2)
		}
		result, err = client.CreateZivaPricingEstimate(ctx, os.Args[2], os.Args[3], os.Args[4])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
