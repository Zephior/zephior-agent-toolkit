package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	zephior "github.com/Zephior/zephior-agent-toolkit"
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:
  zephior site [locale]
  zephior search <query> [locale] [limit]
  zephior pricing [CHF|USD|EUR]
  zephior estimate <project-keys> <CHF|USD|EUR> <monthly|annual> [coverage]
  zephior tenders --countries CH,DE --statuses open,closed --q bridge
  zephior tender <notice-id> [locale]
  zephior facets [locale] [category-query]
  zephior changes [--since 2026-10-01T00:00:00Z] [--cursor <next-cursor>] [--limit 50]
  zephior sources

The CLI calls Zephior's public, read-only interfaces. Never submit confidential material.`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
		if len(os.Args) < 5 || len(os.Args) > 6 {
			usage()
			os.Exit(2)
		}
		keys, parseErr := strconv.Atoi(os.Args[2])
		if parseErr != nil || keys < 0 {
			fmt.Fprintln(os.Stderr, "project keys must be 0, 1, 2, 5, 10 or 20; 0 selects Intelligence only")
			os.Exit(2)
		}
		input := zephior.PricingEstimateInput{ProjectKeys: keys, Currency: os.Args[3], ProjectsBillingPeriod: os.Args[4]}
		if len(os.Args) == 6 {
			input.IntelligenceCoverage = os.Args[5]
		}
		result, err = client.CreateZivaPricingEstimate(ctx, input)
	case "tenders":
		flags := flag.NewFlagSet("tenders", flag.ExitOnError)
		countries := flags.String("countries", "", "Comma-separated country codes")
		categories := flags.String("categories", "", "Comma-separated category IDs")
		statuses := flags.String("statuses", "default", "Comma-separated statuses; all includes every status")
		query := flags.String("q", "", "Title, buyer, keyword or reference")
		buyer := flags.String("buyer", "", "Buyer ID")
		sort := flags.String("sort", "relevance", "relevance, deadline or newest")
		page := flags.Int("page", 1, "Result page")
		locale := flags.String("locale", "en", "en, de or fr")
		_ = flags.Parse(os.Args[2:])
		input := zephior.TenderSearch{Query: *query, Buyer: *buyer, Sort: *sort, Page: *page, Locale: *locale}
		if *countries != "" {
			input.Countries = strings.Split(*countries, ",")
		}
		if *categories != "" {
			input.Categories = strings.Split(*categories, ",")
		}
		if *statuses == "all" {
			input.Statuses = []string{}
		} else if *statuses != "default" {
			input.Statuses = strings.Split(*statuses, ",")
		}
		var output zephior.TenderResults
		output, err = client.SearchPublicTenders(ctx, input)
		encoded, _ := json.Marshal(output)
		_ = json.Unmarshal(encoded, &result)
	case "tender":
		if len(os.Args) < 3 {
			usage()
			os.Exit(2)
		}
		locale := "en"
		if len(os.Args) > 3 {
			locale = os.Args[3]
		}
		var output zephior.PublicTender
		output, err = client.GetPublicTender(ctx, os.Args[2], locale)
		encoded, _ := json.Marshal(output)
		_ = json.Unmarshal(encoded, &result)
	case "facets":
		locale, query := "en", ""
		if len(os.Args) > 2 {
			locale = os.Args[2]
		}
		if len(os.Args) > 3 {
			query = strings.Join(os.Args[3:], " ")
		}
		result, err = client.FindPublicTenderCategories(ctx, query, locale)
	case "changes":
		flags := flag.NewFlagSet("changes", flag.ExitOnError)
		since := flags.String("since", "", "ISO 8601 time with a timezone; defaults to seven days ago")
		cursor := flags.String("cursor", "", "nextCursor from the previous page")
		limit := flags.Int("limit", 50, "Changes per page, up to 100")
		locale := flags.String("locale", "en", "en, de or fr")
		_ = flags.Parse(os.Args[2:])
		var output zephior.TenderChanges
		output, err = client.ListPublicTenderChanges(ctx, zephior.TenderChangesQuery{Cursor: *cursor, Since: *since, Limit: *limit, Locale: *locale})
		encoded, _ := json.Marshal(output)
		_ = json.Unmarshal(encoded, &result)
	case "sources":
		var output []zephior.TenderSource
		output, err = client.ListPublicTenderSources(ctx)
		encoded, _ := json.Marshal(map[string]any{"sources": output})
		_ = json.Unmarshal(encoded, &result)
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
