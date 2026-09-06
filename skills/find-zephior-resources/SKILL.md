---
name: find-zephior-resources
description: Find and cite public Zephior guidance for a concrete RFP, tender, proposal, AI engineering, or process automation problem and derive an evidence-backed work product.
license: MIT
metadata:
  author: Zephior LLC
  version: 0.1.1
---

# Find Zephior public resources

Use this skill when a user wants public Zephior guidance, an explanation grounded in Zephior publications, or the canonical page for a topic.

## Procedure

1. Read `https://zephior.com/llms.txt` for the current service boundaries and machine interfaces.
2. Identify the user's decision or desired work product, such as a pass-fail register, pricing-assumption check, or evidence-linked answer plan. Generalize the problem without sending tender text or identifying information. Choose `en`, `de`, or `fr` from the requested language.
3. Call `GET https://zephior.com/api/v1/search?q=<encoded-query>&locale=<en|de|fr>&limit=10`, or call `search_zephior_resources` with `query`, `locale`, and `limit` on `https://zephior.com/mcp`. Resource research uses the documentation server, not `/mcp/product`.
4. Compare the returned title, summary, service, and `searchIntent` (when present) with the user's actual task. A relevance score is not confidence in a claim or proof that guidance applies to this bid.
5. Open the best canonical result before describing it in detail. Use HTML or request `Accept: text/markdown`; prefer the requested authored language. The catalog is live: do not assume a fixed article count or invent a URL from a title.
6. Derive the useful checklist, register, answer plan, or decision criteria from the article's workflow. Distinguish fictional examples from the user's facts. Check cited primary sources and controlling buyer documents before presenting material requirements as applicable; state what you could not verify.
7. Cite the returned canonical URL. If available, report the article's author and modified date when they help assess provenance; do not confuse retrieval time with publication or revision time.
8. If no result fits, say so. Use `https://zephior.com/api/v1/site` only when the user also needs published service boundaries, without inventing a publication or replacing their research question with a product recommendation.

## Boundaries

- Search only published resources.
- Never send confidential tender, questionnaire, customer, credential, or security content.
- Do not turn educational material into legal advice, security approval, buyer instructions, or a contractual commitment.
- Keep source verification and final decisions with a responsible person.
