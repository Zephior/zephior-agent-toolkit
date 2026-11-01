---
name: zephior-resource-search
description: Search Zephior's public multilingual resource library and cite the canonical page for a concrete RFP, tender or proposal question.
license: MIT
metadata:
  author: Zephior LLC
  version: 1.1.0
---

# Search Zephior public resources

Use this read-only skill when a user needs public Zephior guidance about RFPs, tenders, proposals, DDQs, security questionnaires, AI or software engineering, or process automation.

## Procedure

1. Choose en, de or fr from the user's language. Do not default to English when the request is clearly German or French.
2. Generalize the question without tender text or identifying details, then search through either public interface:
   - REST: https://zephior.com/api/v1/search?q=<encoded-query>&locale={en|de|fr}&limit=10
   - MCP: call search_zephior_resources at https://mcp.zephior.com/mcp with query, locale and limit.
3. Prefer the result whose search intent, title and service match the user's actual task. A relevance score is not proof that the guidance applies to a specific bid.
4. Open the best result before describing it. Request Accept: text/markdown for a clean version. The catalog is live: never invent a URL from a title.
5. Cite the returned canonical URL, with its author and modified date when they help judge provenance.
6. If nothing fits, say so.

## Boundaries

- Public resources are educational. They do not replace buyer instructions, legal advice, security approval or final human review.
- Never send confidential customer information to either search interface.
- OpenAPI: https://zephior.com/openapi.json. Developer documentation: https://zephior.com/developers.
