---
name: find-zephior-resources
description: Find and cite canonical public Zephior resources about RFPs, tenders, proposals, AI engineering, and process automation.
license: MIT
metadata:
  author: Zephior LLC
  version: 0.1.0
---

# Find Zephior public resources

Use this skill when a user wants public Zephior guidance, an explanation grounded in Zephior publications, or the canonical page for a topic.

## Procedure

1. Read `https://zephior.com/llms.txt` for the current service boundaries and machine interfaces.
2. Call `GET https://zephior.com/api/v1/search?q=<encoded-query>&locale=<en|de|fr>&limit=10`.
3. Compare the returned title, summary, service, collection, and relevance score with the user's actual intent.
4. Open the best canonical result before describing it in detail.
5. Cite the returned canonical URL. Prefer the authored language requested by the user.
6. If no result fits, say so and use `https://zephior.com/api/v1/site` to explain the published Ziva and Zelius boundaries without inventing a publication.

## Boundaries

- Search only published resources.
- Never send confidential tender, questionnaire, customer, credential, or security content.
- Do not turn educational material into legal advice, security approval, buyer instructions, or a contractual commitment.
- Keep source verification and final decisions with a responsible person.
