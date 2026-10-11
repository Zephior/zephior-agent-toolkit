# Changes

## 0.4.0

- Published notice search, details and facets in JavaScript, TypeScript, Go and the CLI.
- Multi-country, category and status filters combine with text search and pagination. Omit statuses for active defaults, or send an empty array for every status.
- Pricing estimates now take Intelligence coverage and/or Project Keys. Replace `credits` with `projectKeys` and `billingPeriod` with `projectsBillingPeriod`. Go accepts `PricingEstimateInput`.
- MCP clients initialize their connection, report tool errors and bound response size. Requests have finite deadlines; JavaScript supports caller cancellation and bounded read retries.
- Four portable skills, versioned installation instructions and automated contract checks.

The retired preview, stream and private Tender Universe APIs remain retired. These clients expose published public notices only.
