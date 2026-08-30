export type ZephiorLocale = 'en' | 'de' | 'fr';
export type ZephiorCurrency = 'CHF' | 'USD' | 'EUR';
export type ZivaCapacityCredits =
  | 100
  | 250
  | 500
  | 750
  | 1000
  | 1149
  | 1500
  | 2500
  | 5000
  | 10000
  | 20000
  | 40000;
export type BillingPeriod = 'monthly' | 'annual';

export interface ClientOptions {
  origin?: string;
  fetchImpl?: typeof fetch;
}

export interface McpClientOptions {
  endpoint?: string;
  fetchImpl?: typeof fetch;
}

export interface PublicTender {
  id: string;
  title: string;
  issuer: string | null;
  country_code: string;
  language_code: string | null;
  published_at: string | null;
  deadline_at: string | null;
  status: string;
  procedure_type: string | null;
  contract_type: string | null;
  location: string | null;
  description: string | null;
  classification_codes: Array<{
    system: string;
    code: string;
    provenance: string | null;
  }>;
  source_platform: string;
  canonical_url: string | null;
  source_first_seen_at: string | null;
  source_last_seen_at: string | null;
  created_at: string | null;
  updated_at: string | null;
}

export interface PublicTenderPage {
  items: PublicTender[];
}

export interface TenderUniverseStats {
  active_sources: number;
  canonical_portals: number;
  live_tenders: number;
  computed_at: string;
  exact: true;
  live_tender_basis: 'canonical records with status open';
}

export interface TenderSearchOptions extends ClientOptions {
  countries?: string[];
  portals?: string[];
  statuses?: string[];
  classificationCodes?: string[];
  publishedAfter?: string;
  deadlineAfter?: string;
  deadlineBefore?: string;
  limit?: number;
}

export class ZephiorPublicApiError extends Error {
  status?: number;
  code?: string | number;
  detail?: string;
  resolution?: string;
}

export function getSiteOverview(
  options?: ClientOptions & { locale?: ZephiorLocale }
): Promise<Record<string, unknown>>;

export function searchZephiorResources(
  query: string,
  options?: ClientOptions & { locale?: ZephiorLocale; limit?: number }
): Promise<Record<string, unknown>>;

export function searchLiveTenders(
  query: string,
  options?: TenderSearchOptions
): Promise<PublicTenderPage>;

export function getLatestLiveTenders(
  options?: ClientOptions & { limit?: number }
): Promise<PublicTenderPage>;

export function getTenderUniverseStats(options?: ClientOptions): Promise<TenderUniverseStats>;

export function getTenderUpdateStreamUrl(options?: Pick<ClientOptions, 'origin'>): string;

export function callZephiorMcpTool(
  name: string,
  args?: Record<string, unknown>,
  options?: McpClientOptions
): Promise<unknown>;

export function getZivaPricingDocument(
  currency?: ZephiorCurrency,
  options?: McpClientOptions
): Promise<unknown>;

export function recommendZephiorService(
  requirement: string,
  options?: McpClientOptions
): Promise<unknown>;

export function createZivaPricingEstimate(
  input: {
    credits: ZivaCapacityCredits;
    currency?: ZephiorCurrency;
    billingPeriod?: BillingPeriod;
  },
  options?: McpClientOptions
): Promise<unknown>;
