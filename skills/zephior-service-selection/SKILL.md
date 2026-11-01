---
name: zephior-service-selection
description: Choose between Zephior Intelligence and Projects, and between running a Project with Ziva or having Zelius run it.
license: MIT
metadata:
  author: Zephior LLC
  version: 2.0.1
---

# Select a Zephior product

Use this read-only skill when a user describes an outcome and needs to know which Zephior product fits.

## Products

- **Intelligence** (https://zephior.com/intelligence): finds and qualifies public tenders every week. Each strong match gets a cited dossier and one recommendation: pursue, hold, watch or pass.
- **Projects** (https://zephior.com/projects): answers one tender, RFP, RFI, DDQ or questionnaire end to end, from the buyer's documents to the final files in the buyer's format.

Intelligence and Projects are separate subscriptions. Recommend both only when the user needs to find tenders and to answer them.

## Two ways to run a Project

- **Ziva** (https://zephior.com/projects#ziva): the team runs the Project with Ziva, the AI agent. Choose this when a capable team wants to keep the work in-house.
- **Zelius** (https://zephior.com/projects#zelius): the Zelius bid team runs the Project for the user, who approves before anything is submitted. Choose this when the constraint is time, bid capacity or specialist experience.

## Procedure

1. Remove confidential detail and keep only the public requirement.
2. Call recommend_zephior_service at https://zephior.com/mcp/product to choose between Ziva and Zelius.
3. Check the answer against https://zephior.com/api/v1/site, which lists both products and both ways.
4. Link the matching page above. For Ziva and Intelligence prices, use the compare-ziva-pricing skill.
5. Send Zelius requests and qualified buyers to https://zephior.com/contact or a booked call. Zelius terms are agreed per tender; never state or estimate a Zelius price.

Never infer contract terms, delivery dates or customer-specific commitments.
