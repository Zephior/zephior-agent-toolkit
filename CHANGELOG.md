# Changes

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
