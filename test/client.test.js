import assert from 'node:assert/strict';
import test from 'node:test';
import {
  createZivaPricingEstimate,
  getSiteOverview,
  searchZephiorResources
} from '../src/index.js';

function jsonResponse(payload, init = {}) {
  return new Response(JSON.stringify(payload), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
    ...init
  });
}

test('getSiteOverview requests the selected locale', async () => {
  let requested;
  const payload = await getSiteOverview({
    locale: 'de',
    fetchImpl: async (url) => {
      requested = String(url);
      return jsonResponse({ name: 'Zephior', locale: 'de' });
    }
  });
  assert.equal(requested, 'https://zephior.com/api/v1/site?locale=de');
  assert.equal(payload.locale, 'de');
});

test('searchZephiorResources encodes the public query', async () => {
  let requested;
  await searchZephiorResources('RFP automation', {
    limit: 3,
    fetchImpl: async (url) => {
      requested = new URL(url);
      return jsonResponse({ results: [] });
    }
  });
  assert.equal(requested.searchParams.get('q'), 'RFP automation');
  assert.equal(requested.searchParams.get('limit'), '3');
});

test('createZivaPricingEstimate calls the public product MCP tool', async () => {
  let requestBody;
  const result = await createZivaPricingEstimate(
    { plan: 'Business', currency: 'CHF', billingPeriod: 'monthly' },
    {
      fetchImpl: async (_url, init) => {
        requestBody = JSON.parse(init.body);
        return jsonResponse({
          jsonrpc: '2.0',
          id: 1,
          result: { structuredContent: { amount: 1149, currency: 'CHF' } }
        });
      }
    }
  );
  assert.equal(requestBody.params.name, 'create_ziva_pricing_estimate');
  assert.equal(result.amount, 1149);
});
