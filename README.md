# Zephior Agent Toolkit

[![skills.sh](https://skills.sh/b/Zephior/zephior-agent-toolkit)](https://skills.sh/Zephior/zephior-agent-toolkit)

Official, public agent resources for [Zephior](https://zephior.com): four Agent Skills, an Agent Plugins 1.0 manifest, two remote MCP server declarations, coding-agent rules, and zero-dependency JavaScript and Go clients for published public notices, resources and pricing.

This repository contains no product source code, customer data, private knowledge, credentials, or write-capable integration.

## Install the skills

```bash
npx skills add Zephior/zephior-agent-toolkit
```

The package exposes:

- `find-zephior-resources` for canonical public research and citations;
- `search-public-tenders` for combined country, category, status and text search and source-backed notice details;
- `select-zephior-service` for choosing between Ziva software and Zelius managed bid services;
- `compare-ziva-pricing` for published capacity comparisons and non-binding estimates.

## Use the JavaScript client

Install the tested release directly from GitHub. The npm registry package is not currently published.

```bash
npm install github:Zephior/zephior-agent-toolkit#v0.4.0
```

```js
import {
  createZivaPricingEstimate,
  searchPublicTenders,
  getPublicTender,
  searchZephiorResources,
} from "zephior-agent-toolkit";

const resources = await searchZephiorResources(
  "extract RFP pass fail requirements before drafting",
  {
    locale: "en",
    limit: 5,
  },
);

const estimate = await createZivaPricingEstimate({
  projectKeys: 1,
  currency: "CHF",
  projectsBillingPeriod: "monthly",
});

const tenders = await searchPublicTenders({
  countries: ["CH", "DE"],
  categories: ["construction"],
  statuses: ["open", "closing-soon"],
  q: "bridge",
  locale: "en",
});
if (tenders.items.length) {
  const notice = await getPublicTender(tenders.items[0].id);
  console.log(notice.url, notice.originalSource.url);
}
```

The client uses global Fetch with a 15-second request timeout, caller `signal` support and bounded retries on 429/503. Responses are limited to 8 MiB. MCP calls negotiate initialization and surface tool errors. It has no runtime dependencies.

Filters use OR within selected values and AND across filters and text search. Omit statuses for active defaults; `statuses: []` includes all statuses. Page size is 20. Removing one selected value preserves the remaining array and query. Missing information stays unavailable. Deadlines distinguish exact instants from source calendar dates.

## Research a concrete bid problem

Search for the decision or work product the user needs, not only a product name.
For example, try a pass-fail requirement register, a pricing-assumption check,
or an evidence-linked answer plan. Use `locale: 'de'` or `locale: 'fr'` for a
German or French request.

The toolkit queries the website's published catalog at request time. It does
not bundle article copies or require an update whenever an article is published.
Search results are discovery hints, not sufficient evidence for a detailed answer:

1. Compare the title, summary and `searchIntent` (when present) with the requested
   decision. A relevance score is not a confidence score or proof of applicability.
2. Open the returned canonical URL. Request `Accept: text/markdown` when useful,
   or read the HTML. Keep the requested authored language.
3. Use the article's workflow and inspectable work product; distinguish its
   fictional examples from facts about the user's bid. Check cited primary
   sources and controlling buyer documents for material requirements.
4. Cite the canonical article URL, identify missing evidence, and keep approvals,
   external communications and submission decisions with an authorized person.

Do not paste private tender text or identifying customer information into a
public search query. Generalize the problem first. The installable
[`find-zephior-resources` skill](skills/find-zephior-resources/SKILL.md) provides
the research instructions.

## Use the Go SDK or CLI

```bash
go get github.com/Zephior/zephior-agent-toolkit@latest
go install github.com/Zephior/zephior-agent-toolkit/cmd/zephior@latest

zephior search "extract RFP pass fail requirements before drafting" en 5
zephior pricing CHF
zephior estimate 1 CHF monthly
zephior estimate 0 CHF monthly regional
zephior tenders --countries CH,DE --statuses open,closed --q bridge
zephior facets en
```

The Go API also exposes `SearchPublicTenders`, `GetPublicTender` and `GetPublicTenderFacets`. Pricing uses `PricingEstimateInput{ProjectKeys: 1, Currency: "CHF", ProjectsBillingPeriod: "monthly"}`. Requests have finite timeouts and respect the caller's context cancellation. See [migration notes](CHANGELOG.md).

Select the MCP server explicitly for tool calls:

```go
client := zephior.NewClient()
resources, err := client.CallMCPToolAt(ctx, "/mcp", "search_zephior_resources", map[string]any{
    "query": "extract RFP pass fail requirements before drafting",
    "locale": "en",
    "limit": 5,
})
```

This snippet assumes a `context.Context` named `ctx` and the import
`zephior "github.com/Zephior/zephior-agent-toolkit"`. Handle `err` before using
`resources`.

`CallMCPToolAt` accepts only `/mcp` or `/mcp/product`, resolved against the
client's `BaseURL`. The existing `CallMCPTool` is deprecated but retains its
product-only behavior, so existing callers do not change servers. Typed pricing
methods also remain bound to `/mcp/product`.

## Connect the MCP servers

`mcp.json` declares two public Streamable HTTP servers:

- product decisions: `https://zephior.com/mcp/product`;
- documentation, published tender search and public resources: `https://zephior.com/mcp`.

Both are read-only, require no credential, and expose published information only. The product server also exposes an MCP App UI resource for capacity comparison.

## Canonical contracts

- [Developer documentation](https://zephior.com/developers)
- [llms.txt](https://zephior.com/llms.txt)
- [OpenAPI 3.1](https://zephior.com/openapi.json)
- [MCP discovery](https://zephior.com/.well-known/mcp)
- [Agent Skills discovery](https://zephior.com/.well-known/agent-skills/index.json)
- [Machine-readable pricing](https://zephior.com/pricing.md)

## Safety boundary

Never send customer documents, tender content, questionnaires, credentials, personal data, or other confidential information to the public interfaces. Published material does not replace buyer instructions, legal advice, security approval, pricing authority, source verification, or final human review.

## Development

```bash
npm test
npm pack --dry-run
go test ./...
go vet ./...
node scripts/verify-live.mjs
```

The daily workflow verifies direct Markdown instructions, HTTP and server-rendered
HTML discovery on public and localized pages, live MCP contracts, public tender
queries, pricing and service selection. It also fails before the WebMCP origin-trial
token has fewer than 14 days remaining. Discovery checks require no JavaScript and
reject HTML redirects at `/agents.md`.

Licensed under MIT. Public Zephior content remains subject to the terms linked from the canonical website.
