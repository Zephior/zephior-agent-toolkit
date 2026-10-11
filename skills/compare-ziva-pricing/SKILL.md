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
2. Ask whether Intelligence, Projects or both are needed. Confirm coverage, Project Keys, currency and the applicable Projects billing period.
3. Call `get_ziva_pricing_document` on `https://zephior.com/mcp/product` for a capacity comparison.
4. Call `create_ziva_pricing_estimate` with `currency`, optional `intelligenceCoverage`, optional `projectKeys`, and `projectsBillingPeriod` for Projects. At least one subscription must be selected.
5. Return the published amount without currency conversion, inferred tax, discount, or fee.
6. Label the result “non-binding public estimate” and link `https://zephior.com/pricing` for verification.
7. Treat requirements above 20 Project Keys per month or custom deployment, integration, governance, support, SLA, or managed-delivery requirements as Enterprise scope.
8. Direct Enterprise scope or contractual questions to `https://zephior.com/contact`.

One Project Key unlocks one response project. Intelligence and Projects are independent subscriptions. No customer credits or tokens are billed. A public price is not an accepted order, quote, SLA, or agreement.
