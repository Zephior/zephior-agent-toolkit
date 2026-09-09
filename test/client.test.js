import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
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
    { credits: 1149, currency: 'CHF', billingPeriod: 'monthly' },
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
  assert.deepEqual(requestBody.params.arguments, {
    credits: 1149,
    currency: 'CHF',
    billingPeriod: 'monthly'
  });
  assert.equal(result.amount, 1149);
});

test('package clients do not expose the retired tender machine surface', async () => {
  const files = await Promise.all(
    ['../src/index.js', '../src/index.d.ts', '../client.go', '../cmd/zephior/main.go'].map((path) =>
      readFile(new URL(path, import.meta.url), 'utf8')
    )
  );
  const source = files.join('\n');
  for (const retired of [
    'SearchLiveTenders',
    'GetLatestLiveTenders',
    'GetTenderUniverseStats',
    'TenderUpdateStreamURL',
    '/api/tenders/'
  ]) {
    assert.equal(source.includes(retired), false, `retired public surface remains: ${retired}`);
  }
});
