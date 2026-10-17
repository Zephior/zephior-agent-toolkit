export type ZephiorLocale = "en" | "de" | "fr";
export type ZephiorCurrency = "CHF" | "USD" | "EUR";
export type BillingPeriod = "monthly" | "annual";
export type IntelligenceCoverage =
  "national" | "regional" | "continental" | "global";
export type ProjectKeys = 1 | 2 | 5 | 10 | 20;
export type TenderStatus =
  | "open"
  | "closing-soon"
  | "closed"
  | "awarded"
  | "cancelled"
  | "suspended"
  | "upcoming"
  | "unconfirmed";
export interface ClientOptions {
  origin?: string;
  fetchImpl?: typeof fetch;
  signal?: AbortSignal;
  timeoutMs?: number;
  retries?: number;
}
export interface McpClientOptions extends Omit<ClientOptions, "origin"> {
  endpoint?: string;
}
export class ZephiorPublicApiError extends Error {
  status?: number;
  code?: string | number;
  detail?: string;
  resolution?: string;
}
export interface TenderSearch {
  q?: string;
  countries?: string[];
  categories?: string[];
  statuses?: TenderStatus[];
  buyer?: string;
  sort?: "relevance" | "deadline" | "newest";
  page?: number;
  locale?: ZephiorLocale;
}
export interface PublicTender {
  id: string;
  url: string;
  title: string;
  summary: string;
  locale: ZephiorLocale;
  status: TenderStatus;
  recordedStatus: TenderStatus;
  sourceStatus: string;
  revision: number;
  modifiedAt: string;
  sourceCheckedAt: string;
  stale: boolean;
  country: string;
  deadline: {
    at: string | null;
    date: string | null;
    precision: "instant" | "date" | "unavailable";
    displayTimeZone: string;
    sourceRef: string | null;
  };
  originalSource: { url: string; [key: string]: unknown };
  sources: {
    ref: string;
    url: string | null;
    retrieved_at: string;
    kind: string;
  }[];
  [key: string]: unknown;
}
export interface TenderResults {
  items: PublicTender[];
  total: number;
  page: number;
  limit: 20;
  locale: ZephiorLocale;
  url: string;
  nextPage: number | null;
  filters: TenderSearch;
  filterSemantics: string;
}
export interface TenderFacets {
  countries: { id: string; count: number }[];
  /** At most 100 categories; pass categoryQuery to find others. */
  categories: { id: string; en: string; de: string; fr: string }[];
  categoriesTotal: number;
  statuses: { id: TenderStatus; label: string }[];
  locale: ZephiorLocale;
}
export interface TenderChangesQuery {
  /** nextCursor from the previous page. */
  cursor?: string;
  /** ISO 8601 time with a timezone. Defaults to seven days ago. */
  since?: string;
  limit?: number;
  locale?: ZephiorLocale;
}
export interface TenderChange {
  id: string;
  revision: number;
  changedAt: string;
  visibility: "public" | "withdrawn" | "merged";
  replacementId: string | null;
  url?: string;
  titleEn?: string | null;
  country?: string | null;
  recordedStatus?: string | null;
  deadline?: string | null;
  changedFields?: string[];
}
export interface TenderChanges {
  items: TenderChange[];
  nextCursor: string | null;
  locale: ZephiorLocale;
  semantics: string;
}
export interface TenderSource {
  id: string;
  name: string;
  url: string | null;
  tenders: number;
  currentTenders: number;
  countries: string[];
  lastCheckedAt: string | null;
}
export interface PricingEstimateInput {
  intelligenceCoverage?: IntelligenceCoverage;
  projectKeys?: ProjectKeys;
  currency?: ZephiorCurrency;
  projectsBillingPeriod?: BillingPeriod;
}
export interface PricingEstimate {
  type: "non-binding-public-estimate";
  model: string;
  currency: ZephiorCurrency;
  intelligence: {
    coverage: IntelligenceCoverage;
    amount: number;
    billingPeriod: "monthly";
  } | null;
  projects: {
    keysPerMonth: number;
    amount: number;
    billingPeriod: BillingPeriod;
    keysIssued: number;
    entitlements: string[];
  } | null;
  customPricingRequired: boolean;
  canonicalPricing: string;
  validity: string;
  nextStep: string;
}
export interface ServiceRecommendation {
  requirement: string;
  recommendation: {
    id: "ziva" | "zelius";
    name: string;
    boundary: string;
  } | null;
  resolution: string;
  safety: string;
}
export function getSiteOverview(
  options?: ClientOptions & { locale?: ZephiorLocale },
): Promise<Record<string, unknown>>;
export function searchZephiorResources(
  query: string,
  options?: ClientOptions & { locale?: ZephiorLocale; limit?: number },
): Promise<Record<string, unknown>>;
export function callZephiorMcpTool(
  name: string,
  args?: Record<string, unknown>,
  options?: McpClientOptions,
): Promise<unknown>;
export function getZivaPricingDocument(
  currency?: ZephiorCurrency,
  options?: McpClientOptions,
): Promise<Record<string, unknown>>;
export function recommendZephiorService(
  requirement: string,
  options?: McpClientOptions,
): Promise<ServiceRecommendation>;
export function createZivaPricingEstimate(
  input: PricingEstimateInput,
  options?: McpClientOptions,
): Promise<PricingEstimate>;
export function searchPublicTenders(
  input?: TenderSearch,
  options?: ClientOptions,
): Promise<TenderResults>;
export function getPublicTender(
  id: string,
  options?: ClientOptions & { locale?: ZephiorLocale },
): Promise<PublicTender>;
export function getPublicTenderFacets(
  options?: ClientOptions & { locale?: ZephiorLocale; categoryQuery?: string },
): Promise<TenderFacets>;
export function listPublicTenderChanges(
  input?: TenderChangesQuery,
  options?: ClientOptions,
): Promise<TenderChanges>;
export function listPublicTenderSources(
  options?: ClientOptions,
): Promise<{ sources: TenderSource[] }>;
