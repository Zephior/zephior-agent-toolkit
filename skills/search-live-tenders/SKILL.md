---
name: search-live-tenders
description: Search Zephior's public Tender Universe for current public procurement opportunities and cite the authoritative notice.
license: MIT
metadata:
  author: Zephior LLC
  version: 0.1.0
---

# Search live tenders

Use this skill to find current public procurement notices in Zephior's public
Tender Universe.

## Procedure

1. Keep the query limited to public requirements, industries, locations, and
   classification codes. Never send customer-confidential pursuit strategy,
   credentials, documents, or personal data.
2. Call `POST https://zephior.com/api/tenders/search` with a two-to-500-character
   query, an `open` status filter, and at most 20 results.
3. Use `GET https://zephior.com/api/tenders/latest?limit=10` when the user asks
   for recently published open opportunities without a semantic query.
4. Treat `GET https://zephior.com/api/tenders/stats` as the exact current
   coverage snapshot. Do not substitute stale figures from prose pages.
5. Cite each result's `canonical_url` and identify the issuing authority,
   publication date, deadline, country, and source platform when available.
6. Open the authoritative notice before relying on requirements, eligibility,
   deadlines, documents, or submission instructions.

## Boundaries

- Search results are discovery signals, not a qualification, bid decision, or
  submission instruction.
- Missing fields are unknown. Never infer a deadline, buyer requirement,
  eligibility position, or contract value.
- Keep a responsible human accountable for source verification, qualification,
  legal review, commercial commitments, and submission.
