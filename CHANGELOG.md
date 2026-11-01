# Changes

## 0.6.0

- Skills now match the copies Zephior publishes at `https://zephior.com/.well-known/agent-skills/index.json`, under MIT. Two skills are renamed to their published IDs: `find-zephior-resources` is now `zephior-resource-search` and `select-zephior-service` is now `zephior-service-selection`. Reinstall with `npx skills add Zephior/zephior-agent-toolkit`.
- Product selection follows the current model: Intelligence finds and qualifies tenders (national coverage is free) and Projects answers them, run by your team with Ziva or by the Zelius bid team.
- `search-public-tenders` documents the MCP search filters (sectors, CPV prefixes, place, deadline and value ranges, compact or full view), notice history paging, category paging and the change-feed filters.
- `listPublicTenderChanges` (Go `ListPublicTenderChanges`, CLI `zephior changes`) accepts `until`, `countries`, `sources`, `events` and `sectors`.
- New `browsePublicTenderCategories` (Go `BrowsePublicTenderCategories`, CLI `zephior categories`) pages through every category with `q` and `after`.
- `mcp.json` declares the tenders-only server at `https://mcp.zephior.com/tenders`.
- The MCP contract in `contracts/mcp-tools.json` matches the live servers again, and live verification checks the new filters and that the repository skills match the published ones.

## 0.5.0

- Follow new, amended and withdrawn notices with `listPublicTenderChanges` (Go `ListPublicTenderChanges`, CLI `zephior changes`). Pass `nextCursor` back as `cursor` to continue; the feed starts seven days back by default.
- List the sources behind published notices, with names, official portals, counts and check times, using `listPublicTenderSources` (Go `ListPublicTenderSources`, CLI `zephior sources`).
- Facets now return up to 100 categories with `categoriesTotal`. Pass `categoryQuery` (Go `FindPublicTenderCategories`, CLI `zephior facets en road`) to find other category IDs.
- The documentation MCP server is declared at `https://mcp.zephior.com/mcp`. `https://zephior.com/mcp` keeps working.

## 0.4.0

- Published notice search, details and facets in JavaScript, TypeScript, Go and the CLI.
- Multi-country, category and status filters combine with text search and pagination. Omit statuses for active defaults, or send an empty array for every status.
- Pricing estimates now take Intelligence coverage and/or Project Keys. Replace `credits` with `projectKeys` and `billingPeriod` with `projectsBillingPeriod`. Go accepts `PricingEstimateInput`.
- MCP clients initialize their connection, report tool errors and bound response size. Requests have finite deadlines; JavaScript supports caller cancellation and bounded read retries.
- Four portable skills, versioned installation instructions and automated contract checks.

These clients expose published public notices only.
