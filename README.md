# Zephior Agent Toolkit

[![skills.sh](https://skills.sh/b/tonykzeph/zephior-agent-toolkit)](https://skills.sh/tonykzeph/zephior-agent-toolkit)

Official, public agent resources for [Zephior](https://zephior.com): three Agent Skills, an Agent Plugins 1.0 manifest, two remote MCP server declarations, coding-agent rules, and zero-dependency JavaScript and Go clients for the read-only Discovery API.

This repository contains no product source code, customer data, private knowledge, credentials, or write-capable integration.

## Install the skills

```bash
npx skills add tonykzeph/zephior-agent-toolkit
```

The package exposes:

- `find-zephior-resources` for canonical public research and citations;
- `select-zephior-service` for choosing Ziva, Zelius, Zeke, or Zenith;
- `compare-ziva-pricing` for published plan comparisons and non-binding estimates.

## Use the JavaScript client

```js
import {
  createZivaPricingEstimate,
  searchZephiorResources
} from 'zephior-agent-toolkit';

const resources = await searchZephiorResources('RFP automation', {
  locale: 'en',
  limit: 5
});

const estimate = await createZivaPricingEstimate({
  plan: 'Business',
  currency: 'CHF',
  billingPeriod: 'monthly'
});
```

The client uses the global Fetch API and has no runtime dependencies.

## Use the Go SDK or CLI

```bash
go get github.com/tonykzeph/zephior-agent-toolkit@latest
go install github.com/tonykzeph/zephior-agent-toolkit/cmd/zephior@latest

zephior search "RFP automation" en 5
zephior pricing CHF
zephior estimate Business CHF monthly
```

The Go API exposes `GetSiteOverview`, `SearchResources`, `GetZivaPricingDocument`, `CreateZivaPricingEstimate`, and the generic `CallMCPTool` method.

## Connect the MCP servers

`mcp.json` declares two public Streamable HTTP servers:

- product decisions: `https://zephior.com/mcp/product`;
- documentation and public-resource search: `https://zephior.com/mcp`.

Both are read-only, require no credential, and expose published information only. The product server also exposes an MCP App UI resource for plan comparison.

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
