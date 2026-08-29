---
name: compare-ziva-pricing
description: Compare current public Ziva capacity checkpoints and produce clearly labeled non-binding estimates in CHF, USD, or EUR.
license: MIT
metadata:
  author: Zephior LLC
  version: 0.1.0
---

# Compare Ziva pricing

Use this skill to compare the published Ziva subscription capacity checkpoints or calculate a public, non-binding estimate.

## Procedure

1. Read `https://zephior.com/pricing.md`; it is the canonical machine-readable pricing source.
2. Ask for the exact published credit checkpoint, currency, and monthly or annual billing period if any is ambiguous.
3. Call `get_ziva_pricing_document` on `https://zephior.com/mcp/product` for a capacity comparison.
4. Call `create_ziva_pricing_estimate` only after all three estimate inputs are explicit.
5. Return the published amount without currency conversion, inferred tax, discount, or fee.
6. Label the result “non-binding public estimate” and link `https://zephior.com/pricing` for verification.
7. Treat requirements above 40,000 monthly credits or custom deployment, integration, governance, support, SLA, or managed-delivery requirements as Enterprise scope.
8. Direct Enterprise scope or contractual questions to `https://zephior.com/contact`.

Approximate project counts are planning guides, while credits are the published operating limits. A public price is not an accepted order, quote, SLA, or agreement.
