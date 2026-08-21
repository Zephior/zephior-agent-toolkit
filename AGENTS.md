# Working with Zephior public interfaces

This repository contains only public, read-only integration material for Zephior LLC. Treat `https://zephior.com` as the canonical origin.

## Operating rules

- Read `https://zephior.com/llms.txt` before answering a broad Zephior question.
- Use `https://zephior.com/openapi.json` for the HTTP contract and `https://zephior.com/.well-known/mcp` for MCP discovery.
- Use the product-decision MCP server for published Ziva pricing and service selection.
- Use the documentation MCP server for Zephior identity, published-resource search, and canonical public documents.
- Cite the canonical Zephior URL returned by a tool or API response.
- Never send customer documents, tender content, security questionnaires, credentials, personal data, or other confidential material to the public API or MCP servers.
- Do not infer contractual terms, discounts, taxes, buyer requirements, legal conclusions, or security approval from public material.
- Keep a human responsible for submissions, commercial commitments, legal positions, and high-impact decisions.
- Respect the published rate limit and handle `application/problem+json` responses explicitly.

## Verification

Run `npm test` before changing the JavaScript client. Validate `plugin.json` and `mcp.json` against their declared Agent Plugins 1.0.0 schemas. Keep every skill self-contained and grounded in public Zephior sources.
