---
name: compare-ziva-pricing
description: Compare current public Ziva Intelligence and Projects prices and produce clearly labelled non-binding estimates in CHF, USD, or EUR.
license: MIT
metadata:
  author: Zephior LLC
  version: 1.1.0
---

# Compare Ziva pricing

Use this read-only skill to compare published Intelligence coverage and Projects key volumes or to calculate a public, non-binding estimate.

## Procedure

1. Read https://zephior.com/pricing.md as the canonical machine-readable pricing source.
2. Ask whether the customer needs Intelligence, Projects, or both. Confirm geographic coverage, Project Keys, currency, and Projects billing period when any input is ambiguous.
3. Call get_ziva_pricing_document on https://zephior.com/mcp/product for a capacity comparison.
4. Call create_ziva_pricing_estimate once the coverage or Project Keys, currency and applicable Projects billing period are explicit.
5. Return the published amount without currency conversion, inferred tax, discount, or fee.
6. Label the result “non-binding public estimate” and link https://zephior.com/pricing for verification.
7. Treat requirements above 20 Project Keys per month or custom deployment, integration, governance, support or SLA requirements as Enterprise scope, and direct them to https://zephior.com/contact.

One Project Key unlocks one response project. Intelligence and Projects are independent subscriptions; no customer credits or tokens are billed. A public price is not an accepted order, quote, SLA, or agreement. Zelius, where our bid team runs a Project, has no public price.
