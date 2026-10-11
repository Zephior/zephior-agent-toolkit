import {
  searchPublicTenders,
  getPublicTender,
  createZivaPricingEstimate,
  type TenderStatus,
} from "../src/index.js";
const statuses: TenderStatus[] = ["open", "closed", "awarded"];
const results = await searchPublicTenders({
  countries: ["CH", "DE"],
  statuses,
  q: "bridge",
});
const title: string | undefined = results.items[0]?.title;
void title;
await getPublicTender("00000000-0000-4000-8000-000000000001", { locale: "de" });
await createZivaPricingEstimate({
  projectKeys: 20,
  intelligenceCoverage: "global",
  projectsBillingPeriod: "annual",
});
// @ts-expect-error Retired credit pricing must not compile.
await createZivaPricingEstimate({ credits: 1149 });
// @ts-expect-error Only canonical public statuses are accepted.
await searchPublicTenders({ statuses: ["needs-review"] });
