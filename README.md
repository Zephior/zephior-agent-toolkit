# Zephior Agent Toolkit

[![skills.sh](https://skills.sh/b/Zephior/zephior-agent-toolkit)](https://skills.sh/Zephior/zephior-agent-toolkit)

Official, public agent resources for [Zephior](https://zephior.com): four Agent Skills, an Agent Plugins 1.0 manifest, two remote MCP server declarations, coding-agent rules, and zero-dependency JavaScript and Go clients for the read-only Discovery and Tender Universe APIs.

This repository contains no product source code, customer data, private knowledge, credentials, or write-capable integration.

## Install the skills

```bash
npx skills add Zephior/zephior-agent-toolkit
```

The package exposes:

- `find-zephior-resources` for canonical public research and citations;
- `search-live-tenders` for current public procurement discovery and authoritative-source verification;
- `select-zephior-service` for choosing between Ziva software and Zelius managed bid services;
- `compare-ziva-pricing` for published capacity comparisons and non-binding estimates.

## Use the JavaScript client

```js
import {
  createZivaPricingEstimate,
  getLatestLiveTenders,
  searchLiveTenders,
  searchZephiorResources
} from 'zephior-agent-toolkit';

const resources = await searchZephiorResources('RFP automation', {
  locale: 'en',
  limit: 5
});

const tenders = await searchLiveTenders('cybersecurity services', {
  countries: ['CH'],
  statuses: ['open'],
  limit: 10
});

const latestInEurope = await getLatestLiveTenders({
  region: 'europe',
  limit: 10
});

const estimate = await createZivaPricingEstimate({
  credits: 1149,
  currency: 'CHF',
  billingPeriod: 'monthly'
});
```

The client uses the global Fetch API and has no runtime dependencies.

## Use the Go SDK or CLI

```bash
go get github.com/Zephior/zephior-agent-toolkit@latest
go install github.com/Zephior/zephior-agent-toolkit/cmd/zephior@latest

zephior search "RFP automation" en 5
zephior tenders "cybersecurity services" 10
zephior latest-tenders 10 europe
zephior tender-stats
zephior pricing CHF
zephior estimate 1149 CHF monthly
```

The Go API exposes `GetSiteOverview`, `SearchResources`, `SearchLiveTenders`, `GetLatestLiveTenders`, `GetLatestLiveTendersByRegion`, `GetTenderUniverseStats`, `TenderUpdateStreamURLForRegion`, `GetZivaPricingDocument`, `CreateZivaPricingEstimate`, and the generic `CallMCPTool` method.

## Connect the MCP servers

`mcp.json` declares two public Streamable HTTP servers:

- product decisions: `https://zephior.com/mcp/product`;
- documentation, public-resource search, and bounded Tender Universe tools: `https://zephior.com/mcp`.

Both are read-only, require no credential, and expose published information only. The product server also exposes an MCP App UI resource for capacity comparison.

The website adds a complementary WebMCP layer while a supporting browser has
`/tenders` or the homepage open. WebMCP updates the visible page and is
ephemeral. This toolkit, OpenAPI, and MCP remain the durable interfaces for
off-page agents.

## Canonical contracts

- [Developer documentation](https://zephior.com/developers)
- [llms.txt](https://zephior.com/llms.txt)
- [OpenAPI 3.1](https://zephior.com/openapi.json)
- [MCP discovery](https://zephior.com/.well-known/mcp)
- [Agent Skills discovery](https://zephior.com/.well-known/agent-skills/index.json)
- [Machine-readable pricing](https://zephior.com/pricing.md)
- [Public Tender Universe](https://zephior.com/tenders)
- [Tender search API](https://zephior.com/api/tenders/search)
- [Latest open tenders](https://zephior.com/api/tenders/latest)
- [Exact Tender Universe statistics](https://zephior.com/api/tenders/stats)
- [Compact Tender Universe stream](https://zephior.com/api/tenders/stream?region=global)
- [Hosted live-tender Agent Skill](https://zephior.com/agent-skills/zephior-live-tender-search/SKILL.md)

## Safety boundary

Never send customer documents, tender content, questionnaires, credentials, personal data, or other confidential information to the public interfaces. Published material does not replace buyer instructions, legal advice, security approval, pricing authority, source verification, or final human review.

## Development

```bash
npm test
npm pack --dry-run
go test ./...
```

Licensed under MIT. Public Zephior content remains subject to the terms linked from the canonical website.
