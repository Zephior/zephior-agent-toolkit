---
name: compare-ziva-pricing
description: Compare current public Ziva plans and produce clearly labeled non-binding estimates in CHF, USD, or EUR.
license: MIT
metadata:
  author: Zephior LLC
  version: 0.1.0
---

# Compare Ziva pricing

Use this skill to compare published Ziva subscription plans or calculate a public, non-binding estimate.

## Procedure

1. Read `https://zephior.com/pricing.md`; it is the canonical machine-readable pricing source.
2. Ask for the plan, published currency, and monthly or annual billing period if any is ambiguous.
3. Call `get_ziva_pricing_document` on `https://zephior.com/mcp/product` for a plan comparison.
4. Call `create_ziva_pricing_estimate` only after all three estimate inputs are explicit.
5. Return the published amount without currency conversion, inferred tax, discount, or fee.
6. Label the result “non-binding public estimate” and link `https://zephior.com/pricing` for verification.
7. Treat Enterprise as custom pricing; do not convert `null` or `Custom` into zero.
8. Direct custom scope or contractual questions to `https://zephior.com/contact`.

Approximate project counts are planning guides, while credits are the published operating limits. A public price is not an accepted order, quote, SLA, or agreement.
