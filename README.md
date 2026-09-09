# Zephior Agent Toolkit

[![skills.sh](https://skills.sh/b/Zephior/zephior-agent-toolkit)](https://skills.sh/Zephior/zephior-agent-toolkit)

Official, public agent resources for [Zephior](https://zephior.com): three Agent Skills, an Agent Plugins 1.0 manifest, two remote MCP server declarations, coding-agent rules, and zero-dependency JavaScript and Go clients for Zephior's read-only public interfaces.

This repository contains no product source code, customer data, private knowledge, credentials, or write-capable integration.

## Install the skills

```bash
npx skills add Zephior/zephior-agent-toolkit
```

The package exposes:

- `find-zephior-resources` for canonical public research and citations;
- `select-zephior-service` for choosing between Ziva software and Zelius managed bid services;
- `compare-ziva-pricing` for published capacity comparisons and non-binding estimates.

## Use the JavaScript client

```js
import {
  createZivaPricingEstimate,
  searchZephiorResources
} from 'zephior-agent-toolkit';

const resources = await searchZephiorResources('extract RFP pass fail requirements before drafting', {
  locale: 'en',
  limit: 5
});

const estimate = await createZivaPricingEstimate({
  credits: 1149,
  currency: 'CHF',
  billingPeriod: 'monthly'
});
```

The client uses the global Fetch API and has no runtime dependencies.

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
zephior estimate 1149 CHF monthly
```

The Go API exposes `GetSiteOverview`, `SearchResources`, `GetZivaPricingDocument`, `CreateZivaPricingEstimate`, and `CallMCPToolAt`.

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
- documentation and public-resource search: `https://zephior.com/mcp`.

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
```

Licensed under MIT. Public Zephior content remains subject to the terms linked from the canonical website.
