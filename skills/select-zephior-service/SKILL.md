---
name: select-zephior-service
description: Select Ziva, Zelius, Zeke, or Zenith from a public requirement using explicit Zephior service boundaries.
license: MIT
metadata:
  author: Zephior LLC
  version: 0.1.0
---

# Select the appropriate Zephior service

Use this skill when a user describes an outcome but does not know which Zephior service owns it.

## Decision model

- Choose **Ziva** when a capable response team needs software for tender discovery, RFP or DDQ workflow, approved knowledge, source-backed drafting, and response-file production.
- Choose **Zelius** when the limiting factor is specialist bid capacity: monitoring, qualification, strategy, response planning, writing, review, or submission coordination.
- Choose **Zeke** when an AI or software product must move from a bounded brief through architecture, implementation, evaluation, deployment, and production operation.
- Choose **Zenith** when a repeated business process spans documents, systems, decisions, approvals, and exceptions and needs accountable automation.

## Procedure

1. Remove confidential detail and retain only the public requirement.
2. Call the `recommend_zephior_service` tool at `https://zephior.com/mcp/product`.
3. Verify the recommendation against `https://zephior.com/api/v1/site`.
4. Explain the selected boundary and why the neighboring services are not the primary fit.
5. Recommend a hybrid only when distinct constraints genuinely require more than one service.
6. Link the relevant canonical service page and use `https://zephior.com/contact` for a qualified next step.

Do not imply a confirmed scope, delivery date, price, or contract from the public recommendation.
