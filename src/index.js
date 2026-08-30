const DEFAULT_ORIGIN = 'https://zephior.com';

export class ZephiorPublicApiError extends Error {
  constructor(message, { status, code, detail, resolution } = {}) {
    super(message);
    this.name = 'ZephiorPublicApiError';
    this.status = status;
    this.code = code;
    this.detail = detail;
    this.resolution = resolution;
  }
}

function normalizedOrigin(origin) {
  return String(origin || DEFAULT_ORIGIN).replace(/\/$/, '');
}

function selectedFetch(fetchImpl) {
  const implementation = fetchImpl || globalThis.fetch;
  if (typeof implementation !== 'function') {
    throw new TypeError('A Fetch API implementation is required.');
  }
  return implementation;
}

async function readJson(response) {
  const payload = await response.json().catch(() => ({}));
  if (!response.ok || payload.error) {
    const problem = payload.error || payload;
    throw new ZephiorPublicApiError(
      problem.title || problem.message || `Zephior request failed with HTTP ${response.status}`,
      {
        status: response.status,
        code: problem.code,
        detail: problem.detail,
        resolution: problem.resolution
      }
    );
  }
  return payload;
}

export async function getSiteOverview({ locale = 'en', origin = DEFAULT_ORIGIN, fetchImpl } = {}) {
  const url = new URL('/api/v1/site', normalizedOrigin(origin));
  url.searchParams.set('locale', locale);
  const response = await selectedFetch(fetchImpl)(url, {
    headers: { Accept: 'application/json' }
  });
  return readJson(response);
}

export async function searchZephiorResources(
  query,
  { locale = 'en', limit = 5, origin = DEFAULT_ORIGIN, fetchImpl } = {}
) {
  if (typeof query !== 'string' || query.trim().length < 2) {
    throw new TypeError('query must contain at least two characters');
  }
  const url = new URL('/api/v1/search', normalizedOrigin(origin));
  url.searchParams.set('q', query.trim());
  url.searchParams.set('locale', locale);
  url.searchParams.set('limit', String(limit));
  const response = await selectedFetch(fetchImpl)(url, {
    headers: { Accept: 'application/json' }
  });
  return readJson(response);
}

function boundedTenderLimit(limit, fallback) {
  const value = Number(limit ?? fallback);
  if (!Number.isInteger(value) || value < 1 || value > 20) {
    throw new TypeError('limit must be an integer between 1 and 20');
  }
  return value;
}

const TENDER_REGION_GROUPS = new Set([
  'europe',
  'mena',
  'africa',
  'north-america',
  'latin-america',
  'central-asia',
  'south-asia',
  'east-asia',
  'southeast-asia',
  'oceania'
]);

function boundedTenderRegion(region, { allowAuto = true } = {}) {
  const value = String(region ?? 'global').trim().toLowerCase();
  if (value === 'global' || (allowAuto && value === 'auto') || TENDER_REGION_GROUPS.has(value)) {
    return value;
  }
  if (/^[a-z]{2}$/.test(value)) return value.toUpperCase();
  throw new TypeError(
    allowAuto
      ? 'region must be global, auto, a regional slug, or an ISO alpha-2 country code'
      : 'region must be global, a regional slug, or an ISO alpha-2 country code'
  );
}

export async function searchLiveTenders(
  query,
  {
    countries = [],
    portals = [],
    statuses = ['open'],
    classificationCodes = [],
    publishedAfter,
    deadlineAfter,
    deadlineBefore,
    limit = 20,
    origin = DEFAULT_ORIGIN,
    fetchImpl
  } = {}
) {
  if (typeof query !== 'string' || query.trim().length < 2 || query.trim().length > 500) {
    throw new TypeError('query must contain between two and 500 characters');
  }
  const response = await selectedFetch(fetchImpl)(
    new URL('/api/tenders/search', normalizedOrigin(origin)),
    {
      method: 'POST',
      headers: { Accept: 'application/json', 'Content-Type': 'application/json' },
      body: JSON.stringify({
        query: query.trim(),
        countries,
        portals,
        statuses,
        classification_codes: classificationCodes,
        published_after: publishedAfter,
        deadline_after: deadlineAfter,
        deadline_before: deadlineBefore,
        limit: boundedTenderLimit(limit, 20)
      })
    }
  );
  return readJson(response);
}

export async function getLatestLiveTenders(
  { limit = 10, region = 'global', origin = DEFAULT_ORIGIN, fetchImpl } = {}
) {
  const url = new URL('/api/tenders/latest', normalizedOrigin(origin));
  url.searchParams.set('limit', String(boundedTenderLimit(limit, 10)));
  url.searchParams.set('region', boundedTenderRegion(region));
  const response = await selectedFetch(fetchImpl)(url, {
    headers: { Accept: 'application/json' }
  });
  return readJson(response);
}

export async function getTenderUniverseStats({ origin = DEFAULT_ORIGIN, fetchImpl } = {}) {
  const response = await selectedFetch(fetchImpl)(
    new URL('/api/tenders/stats', normalizedOrigin(origin)),
    { headers: { Accept: 'application/json' } }
  );
  return readJson(response);
}

export function getTenderUpdateStreamUrl({ region = 'global', origin = DEFAULT_ORIGIN } = {}) {
  const url = new URL('/api/tenders/stream', normalizedOrigin(origin));
  url.searchParams.set('region', boundedTenderRegion(region, { allowAuto: false }));
  return url.toString();
}

export async function callZephiorMcpTool(
  name,
  args = {},
  { endpoint = `${DEFAULT_ORIGIN}/mcp/product`, fetchImpl } = {}
) {
  if (typeof name !== 'string' || !name) throw new TypeError('name is required');
  const response = await selectedFetch(fetchImpl)(endpoint, {
    method: 'POST',
    headers: {
      Accept: 'application/json, text/event-stream',
      'Content-Type': 'application/json',
      'MCP-Protocol-Version': '2025-06-18'
    },
    body: JSON.stringify({
      jsonrpc: '2.0',
      id: 1,
      method: 'tools/call',
      params: { name, arguments: args }
    })
  });
  const payload = await readJson(response);
  if (payload.error) {
    throw new ZephiorPublicApiError(payload.error.message || 'MCP tool call failed', {
      status: response.status,
      code: payload.error.code,
      detail: payload.error.data?.message,
      resolution: payload.error.data?.resolution
    });
  }
  return payload.result?.structuredContent ?? payload.result;
}

export function getZivaPricingDocument(currency = 'CHF', options = {}) {
  return callZephiorMcpTool('get_ziva_pricing_document', { currency }, options);
}

export function recommendZephiorService(requirement, options = {}) {
  return callZephiorMcpTool('recommend_zephior_service', { requirement }, options);
}

export function createZivaPricingEstimate(
  { credits, currency = 'CHF', billingPeriod = 'monthly' },
  options = {}
) {
  return callZephiorMcpTool(
    'create_ziva_pricing_estimate',
    { credits, currency, billingPeriod },
    options
  );
}
