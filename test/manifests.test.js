import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";
import Ajv from "ajv/dist/2020.js";
for (const name of ["plugin", "mcp"])
  test(`${name} manifest matches its declared public schema`, async () => {
    const manifest = JSON.parse(
      await readFile(new URL(`../${name}.json`, import.meta.url), "utf8"),
    );
    const response = await fetch(manifest.$schema, {
      signal: AbortSignal.timeout(15000),
    });
    assert.equal(response.status, 200);
    const schema = await response.json();
    const ajv = new Ajv({ strict: false, validateFormats: false });
    const validate = ajv.compile(schema);
    assert.equal(validate(manifest), true, JSON.stringify(validate.errors));
    if (name === "plugin")
      assert.equal(
        manifest.version,
        JSON.parse(
          await readFile(new URL("../package.json", import.meta.url), "utf8"),
        ).version,
      );
  });
