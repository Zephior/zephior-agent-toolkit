---
name: search-live-tenders
description: Search Zephior's public Tender Universe for current public procurement opportunities and cite the authoritative notice.
license: MIT
metadata:
  author: Zephior LLC
  version: 0.1.2
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
3. Use `GET https://zephior.com/api/tenders/latest?region=global&limit=10` when the user asks
   for recently published open opportunities with a non-future publication
   timestamp, a future submission deadline, and no semantic query. Set region to
   a documented regional slug or ISO alpha-2 country code when requested. Read
   `scope.requested`, `scope.effective`, and `scope.broadened` before describing
   the geographic scope or subscribing to the effective stream.
4. Treat `GET https://zephior.com/api/tenders/stats` as the exact current
   coverage snapshot. Do not substitute stale figures from prose pages.
5. Cite each result's `canonical_url` and identify the issuing authority,
   publication date, deadline, country, and source platform when available.
6. Open the authoritative notice before relying on requirements, eligibility,
   deadlines, documents, or submission instructions.

## Agent-native interfaces

- Hosted Skill: `https://zephior.com/agent-skills/zephior-live-tender-search/SKILL.md`
- OpenAPI: `https://zephior.com/openapi.json`
- MCP tools: `search_public_tenders`, `list_latest_open_tenders`, and
  `get_tender_universe_stats` at `https://zephior.com/mcp`
- WebMCP: supporting agent browsers discover visible tender tools while
  `https://zephior.com/tenders` or the relevant homepage surface is open.

## Boundaries

- Search results are discovery signals, not a qualification, bid decision, or
  submission instruction.
- Missing fields are unknown. Never infer a deadline, buyer requirement,
  eligibility position, or contract value.
- Keep a responsible human accountable for source verification, qualification,
  legal review, commercial commitments, and submission.
