const DEFAULT_ORIGIN = "https://zephior.com";
async function request(url, init, options = {}) {
  const signal = options.signal
    ? AbortSignal.any([
        options.signal,
        AbortSignal.timeout(options.timeoutMs ?? 15000),
      ])
    : AbortSignal.timeout(options.timeoutMs ?? 15000);
  const retries = Math.max(0, Math.min(2, options.retries ?? 1));
  for (let attempt = 0; ; attempt++) {
    signal.throwIfAborted();
    const response = await selectedFetch(options.fetchImpl)(url, {
      ...init,
      signal,
      redirect: "error",
    });
    if (![429, 503].includes(response.status) || attempt >= retries)
      return response;
    const retryAfter = response.headers.get("Retry-After");
    const delay =
      retryAfter && /^\d+$/.test(retryAfter) ? Number(retryAfter) * 1000 : 500;
    if (delay > 5000) return response;
    await response.body?.cancel();
    await new Promise((resolve, reject) => {
      const abort = () => {
        clearTimeout(timer);
        reject(signal.reason);
      };
      const timer = setTimeout(() => {
        signal.removeEventListener("abort", abort);
        resolve();
      }, delay);
      signal.addEventListener("abort", abort, { once: true });
    });
  }
}

export class ZephiorPublicApiError extends Error {
  constructor(message, { status, code, detail, resolution } = {}) {
    super(message);
    this.name = "ZephiorPublicApiError";
    this.status = status;
    this.code = code;
    this.detail = detail;
    this.resolution = resolution;
  }
}

function normalizedOrigin(origin) {
  return String(origin || DEFAULT_ORIGIN).replace(/\/$/, "");
}

function selectedFetch(fetchImpl) {
  const implementation = fetchImpl || globalThis.fetch;
  if (typeof implementation !== "function") {
    throw new TypeError("A Fetch API implementation is required.");
  }
  return implementation;
}

async function readJson(response) {
  const maximum = 8 * 1024 * 1024;
  if (Number(response.headers.get("Content-Length")) > maximum)
    throw new ZephiorPublicApiError("Response exceeds 8 MiB");
  const reader = response.body?.getReader();
  let text = "";
  if (reader) {
    const decoder = new TextDecoder();
    let bytes = 0;
    try {
      while (true) {
        const part = await reader.read();
        if (part.done) break;
        bytes += part.value.byteLength;
        if (bytes > maximum) {
          await reader.cancel();
          throw new ZephiorPublicApiError("Response exceeds 8 MiB");
        }
        text += decoder.decode(part.value, { stream: true });
      }
      text += decoder.decode();
    } finally {
      reader.releaseLock();
    }
  }
  if (response.headers.get("content-type")?.includes("text/event-stream")) {
    const event = text
      .split(/\r?\n\r?\n/)
      .map((block) =>
        block
          .split(/\r?\n/)
          .filter((line) => line.startsWith("data:"))
          .map((line) => line.slice(5).trimStart())
          .join("\n"),
      )
      .find((data) => data && JSON.parse(data).id !== undefined);
    text = event ?? "{}";
  }
  const payload = JSON.parse(text || "{}");
  if (!response.ok || payload.error) {
    const problem = payload.error || payload;
    throw new ZephiorPublicApiError(
      problem.title ||
        problem.message ||
        `Zephior request failed with HTTP ${response.status}`,
      {
        status: response.status,
        code: problem.code,
        detail: problem.detail,
        resolution: problem.resolution ?? problem.data?.resolution,
      },
    );
  }
  return payload;
}

export async function getSiteOverview({
  locale = "en",
  origin = DEFAULT_ORIGIN,
  ...options
} = {}) {
  const url = new URL("/api/v1/site", normalizedOrigin(origin));
  url.searchParams.set("locale", locale);
  const response = await request(
    url,
    { headers: { Accept: "application/json" } },
    options,
  );
  return readJson(response);
}

export async function searchZephiorResources(
  query,
  { locale = "en", limit = 5, origin = DEFAULT_ORIGIN, ...options } = {},
) {
  if (typeof query !== "string" || query.trim().length < 2) {
    throw new TypeError("query must contain at least two characters");
  }
  const url = new URL("/api/v1/search", normalizedOrigin(origin));
  url.searchParams.set("q", query.trim());
  url.searchParams.set("locale", locale);
  url.searchParams.set("limit", String(limit));
  const response = await request(
    url,
    { headers: { Accept: "application/json" } },
    options,
  );
  return readJson(response);
}

export async function callZephiorMcpTool(
  name,
  args = {},
  { endpoint = `${DEFAULT_ORIGIN}/mcp/product`, ...options } = {},
) {
  if (typeof name !== "string" || !name)
    throw new TypeError("name is required");
  const initialize = await request(
    endpoint,
    {
      method: "POST",
      headers: {
        Accept: "application/json, text/event-stream",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        jsonrpc: "2.0",
        id: 0,
        method: "initialize",
        params: {
          protocolVersion: "2025-06-18",
          capabilities: {},
          clientInfo: { name: "zephior-agent-toolkit", version: "0.4.0" },
        },
      }),
    },
    options,
  );
  const handshake = await readJson(initialize);
  const protocolVersion = handshake.result?.protocolVersion;
  if (protocolVersion !== "2025-06-18")
    throw new ZephiorPublicApiError("Unsupported MCP protocol negotiated");
  const initialized = await request(
    endpoint,
    {
      method: "POST",
      headers: {
        Accept: "application/json, text/event-stream",
        "Content-Type": "application/json",
        "MCP-Protocol-Version": protocolVersion,
      },
      body: JSON.stringify({
        jsonrpc: "2.0",
        method: "notifications/initialized",
      }),
    },
    options,
  );
  if (initialized.status !== 202)
    throw new ZephiorPublicApiError(
      "MCP initialization notification was rejected",
    );
  const response = await request(
    endpoint,
    {
      method: "POST",
      headers: {
        Accept: "application/json, text/event-stream",
        "Content-Type": "application/json",
        "MCP-Protocol-Version": "2025-06-18",
      },
      body: JSON.stringify({
        jsonrpc: "2.0",
        id: 1,
        method: "tools/call",
        params: { name, arguments: args },
      }),
    },
    options,
  );
  const payload = await readJson(response);
  if (payload.error) {
    throw new ZephiorPublicApiError(
      payload.error.message || "MCP tool call failed",
      {
        status: response.status,
        code: payload.error.code,
        detail: payload.error.data?.message,
        resolution: payload.error.data?.resolution,
      },
    );
  }
  if (payload.result?.isError)
    throw new ZephiorPublicApiError(
      payload.result.content?.find((item) => item.type === "text")?.text ||
        "Public tool could not complete",
    );
  return payload.result?.structuredContent ?? payload.result;
}

export function getZivaPricingDocument(currency = "CHF", options = {}) {
  return callZephiorMcpTool("get_ziva_pricing_document", { currency }, options);
}

export function recommendZephiorService(requirement, options = {}) {
  return callZephiorMcpTool(
    "recommend_zephior_service",
    { requirement },
    options,
  );
}

export function createZivaPricingEstimate(
  {
    intelligenceCoverage,
    projectKeys,
    currency = "CHF",
    projectsBillingPeriod = "monthly",
  },
  options = {},
) {
  if (!intelligenceCoverage && !projectKeys)
    throw new TypeError("Select intelligenceCoverage, projectKeys, or both");
  return callZephiorMcpTool(
    "create_ziva_pricing_estimate",
    {
      currency,
      ...(intelligenceCoverage ? { intelligenceCoverage } : {}),
      ...(projectKeys ? { projectKeys, projectsBillingPeriod } : {}),
    },
    options,
  );
}
export async function searchPublicTenders(
  input = {},
  { origin = DEFAULT_ORIGIN, ...options } = {},
) {
  const url = new URL("/api/v1/tenders", normalizedOrigin(origin));
  for (const [key, value] of Object.entries(input))
    if (value !== undefined)
      url.searchParams.set(
        key,
        Array.isArray(value) ? value.join(",") : String(value),
      );
  return readJson(
    await request(url, { headers: { Accept: "application/json" } }, options),
  );
}
export async function getPublicTender(
  id,
  { locale = "en", origin = DEFAULT_ORIGIN, ...options } = {},
) {
  if (
    !/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(id)
  )
    throw new TypeError("A public notice UUID is required");
  const url = new URL(`/api/v1/tenders/${id}`, normalizedOrigin(origin));
  url.searchParams.set("locale", locale);
  return readJson(
    await request(url, { headers: { Accept: "application/json" } }, options),
  );
}
export async function getPublicTenderFacets({
  locale = "en",
  origin = DEFAULT_ORIGIN,
  ...options
} = {}) {
  const url = new URL("/api/v1/tenders/facets", normalizedOrigin(origin));
  url.searchParams.set("locale", locale);
  return readJson(
    await request(url, { headers: { Accept: "application/json" } }, options),
  );
}
