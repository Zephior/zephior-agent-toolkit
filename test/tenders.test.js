import assert from "node:assert/strict";
import test from "node:test";
import {
  searchPublicTenders,
  getPublicTender,
  getPublicTenderFacets,
  callZephiorMcpTool,
  ZephiorPublicApiError,
} from "../src/index.js";
const json = (data, init = {}) =>
  new Response(JSON.stringify(data), {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
test("multi-select encodes all statuses and preserves search and pagination", async () => {
  let url;
  await searchPublicTenders(
    {
      countries: ["CH", "DE"],
      categories: ["construction", "it"],
      statuses: [],
      q: "bridge ' OR 1=1 --",
      page: 2,
    },
    {
      fetchImpl: async (input) => {
        url = new URL(input);
        return json({ items: [] });
      },
    },
  );
  assert.equal(url.searchParams.get("countries"), "CH,DE");
  assert.equal(url.searchParams.get("statuses"), "");
  assert.equal(url.searchParams.get("page"), "2");
  assert.equal(url.searchParams.get("q"), "bridge ' OR 1=1 --");
  await searchPublicTenders(
    {},
    {
      fetchImpl: async (input) => {
        assert.equal(new URL(input).searchParams.has("statuses"), false);
        return json({ items: [] });
      },
    },
  );
});
test("details reject path injection before fetching and facets preserve locale", async () => {
  await assert.rejects(getPublicTender("../private"), TypeError);
  await getPublicTenderFacets({
    locale: "fr",
    fetchImpl: async (input) => {
      assert.equal(new URL(input).searchParams.get("locale"), "fr");
      return json({ statuses: [] });
    },
  });
});
test("MCP negotiates initialization, accepts SSE, and reports execution failures", async () => {
  const calls = [];
  const fetchImpl = async (_, init) => {
    const body = JSON.parse(init.body);
    calls.push(body.method);
    if (body.method === "initialize")
      return json({ result: { protocolVersion: "2025-06-18" } });
    if (body.method === "notifications/initialized")
      return new Response(null, { status: 202 });
    return new Response(
      "event: message\ndata: " +
        JSON.stringify({
          jsonrpc: "2.0",
          id: 1,
          result: {
            isError: true,
            content: [{ type: "text", text: "Source unavailable" }],
          },
        }) +
        "\n\n",
      { headers: { "Content-Type": "text/event-stream" } },
    );
  };
  await assert.rejects(
    callZephiorMcpTool("get_public_tender", {}, { fetchImpl }),
    /Source unavailable/,
  );
  assert.deepEqual(calls, [
    "initialize",
    "notifications/initialized",
    "tools/call",
  ]);
});
test("cancellation and advertised oversized responses are bounded", async () => {
  const controller = new AbortController();
  controller.abort();
  let calls = 0;
  await assert.rejects(
    searchPublicTenders(
      {},
      {
        signal: controller.signal,
        fetchImpl: async () => {
          calls++;
          return json({});
        },
      },
    ),
  );
  assert.equal(calls, 0);
  await assert.rejects(
    searchPublicTenders(
      {},
      {
        fetchImpl: async () =>
          new Response("{}", { headers: { "Content-Length": "9000000" } }),
      },
    ),
    ZephiorPublicApiError,
  );
});
