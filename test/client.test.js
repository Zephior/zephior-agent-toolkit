import assert from 'node:assert/strict';
import test from 'node:test';
import {
  createZivaPricingEstimate,
  getLatestLiveTenders,
  getSiteOverview,
  getTenderUniverseStats,
  getTenderUpdateStreamUrl,
  searchLiveTenders,
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

test('searchLiveTenders sends bounded public filters', async () => {
  let requestBody;
  const result = await searchLiveTenders('cybersecurity services', {
    countries: ['CH'],
    limit: 5,
    fetchImpl: async (_url, init) => {
      requestBody = JSON.parse(init.body);
      return jsonResponse({ items: [] });
    }
  });

  assert.deepEqual(requestBody, {
    query: 'cybersecurity services',
    countries: ['CH'],
    portals: [],
    statuses: ['open'],
    classification_codes: [],
    limit: 5
  });
  assert.deepEqual(result, { items: [] });
});

test('latest, stats and stream use the published tender endpoints', async () => {
  const requested = [];
  const fetchImpl = async (url) => {
    requested.push(String(url));
    return jsonResponse(String(url).includes('/stats') ? { exact: true } : { items: [] });
  };

  await getLatestLiveTenders({ limit: 3, fetchImpl });
  await getTenderUniverseStats({ fetchImpl });

  assert.deepEqual(requested, [
    'https://zephior.com/api/tenders/latest?limit=3',
    'https://zephior.com/api/tenders/stats'
  ]);
  assert.equal(getTenderUpdateStreamUrl(), 'https://zephior.com/api/tenders/stream');
});
