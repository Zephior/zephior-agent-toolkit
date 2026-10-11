import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import Ajv from "ajv/dist/2020.js";
const contracts = JSON.parse(
  await readFile(
    new URL("../contracts/mcp-tools.json", import.meta.url),
    "utf8",
  ),
);
const ajv = new Ajv({ strict: false, validateFormats: false });
import {
  searchPublicTenders,
  getPublicTender,
  getPublicTenderFacets,
  createZivaPricingEstimate,
  recommendZephiorService,
} from "../src/index.js";
const origin = process.argv[2] ?? "https://zephior.com";
const options = { origin, retries: 2 };
const endpoint = `${origin}/mcp/product`;
for (const [path, tools] of Object.entries(contracts)) {
  const response = await fetch(origin + path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ jsonrpc: "2.0", id: 1, method: "tools/list" }),
  });
  assert.equal(response.status, 200);
  const live = (await response.json()).result.tools;
  const cardPath =
    path === "/mcp"
      ? "/.well-known/mcp/docs-server-card.json"
      : "/.well-known/mcp/server-card.json";
  const card = await (await fetch(origin + cardPath)).json();
  const schema = await (
    await fetch(origin + new URL(card.$schema).pathname)
  ).json();
  const validateCard = ajv.getSchema(schema.$id) ?? ajv.compile(schema);
  assert.ok(validateCard(card), JSON.stringify(validateCard.errors));
  assert.deepEqual(
    card.tools.map((tool) => tool.name).sort(),
    live.map((tool) => tool.name).sort(),
  );
  const initialization = await fetch(origin + path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      jsonrpc: "2.0",
      id: 2,
      method: "initialize",
      params: {
        protocolVersion: "2025-06-18",
        capabilities: {},
        clientInfo: { name: "zephior-toolkit-audit", version: "0.4.0" },
      },
    }),
  });
  const server = (await initialization.json()).result.serverInfo;
  assert.equal(card.name, server.name);
  assert.equal(card.version, server.version);
  for (const [name, expected] of Object.entries(tools)) {
    const actual = live.find((tool) => tool.name === name);
    assert.ok(actual, `Missing documented tool ${name}`);
    assert.deepEqual(
      actual.inputSchema,
      expected.inputSchema,
      `Input contract drift: ${name}`,
    );
    assert.deepEqual(
      actual.outputSchema ?? null,
      expected.outputSchema,
      `Output contract drift: ${name}`,
    );
  }
}
const validateOutput = (path, name, value) => {
  const schema = contracts[path][name].outputSchema;
  if (schema) {
    const validate = ajv.compile(schema);
    assert.ok(validate(value), `${name}: ${JSON.stringify(validate.errors)}`);
  }
};
const results = await searchPublicTenders(
  { countries: ["CH", "DE"], statuses: [], locale: "en" },
  options,
);
validateOutput("/mcp", "search_public_tenders", results);
assert.ok(Array.isArray(results.items));
assert.equal(results.limit, 20);
assert.ok(results.items.every((item) => ["CH", "DE"].includes(item.country)));
if (results.items.length)
  assert.equal(
    (await getPublicTender(results.items[0].id, options)).id,
    results.items[0].id,
  );
assert.equal((await getPublicTenderFacets(options)).statuses.length, 8);
const estimate = await createZivaPricingEstimate(
  { projectKeys: 1, currency: "CHF", projectsBillingPeriod: "monthly" },
  { endpoint },
);
validateOutput("/mcp/product", "create_ziva_pricing_estimate", estimate);
assert.equal(estimate.type, "non-binding-public-estimate");
assert.equal(estimate.projects.keysPerMonth, 1);
assert.equal(
  (await recommendZephiorService("managed tender monitoring", { endpoint }))
    .recommendation.id,
  "zelius",
);
const openapi = await (await fetch(`${origin}/openapi.json`)).json();
assert.ok(openapi.paths["/tenders"]);
assert.ok(openapi.paths["/tenders/{id}"]);
const html = await fetch(`${origin}/tenders`, {
  headers: { Accept: "text/html" },
});
if (origin === "https://zephior.com") {
  const token = html.headers.get("Origin-Trial");
  assert.ok(token, "Missing WebMCP origin-trial token");
  const trial = JSON.parse(
    Buffer.from(token, "base64").subarray(69).toString(),
  );
  assert.equal(trial.feature, "WebMCP");
  assert.equal(trial.origin, "https://zephior.com:443");
  assert.ok(
    trial.expiry * 1000 > Date.now() + 14 * 86400000,
    "Renew the WebMCP trial token before its remaining lifetime falls below 14 days",
  );
}
console.log("Live toolkit contract verification passed");
