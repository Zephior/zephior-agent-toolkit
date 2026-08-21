export type ZephiorLocale = 'en' | 'de' | 'fr';
export type ZephiorCurrency = 'CHF' | 'USD' | 'EUR';
export type ZivaPlan = 'Starter' | 'Business' | 'Advanced' | 'Enterprise';
export type BillingPeriod = 'monthly' | 'annual';

export interface ClientOptions {
  origin?: string;
  fetchImpl?: typeof fetch;
}

export interface McpClientOptions {
  endpoint?: string;
  fetchImpl?: typeof fetch;
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
  input: { plan: ZivaPlan; currency?: ZephiorCurrency; billingPeriod?: BillingPeriod },
  options?: McpClientOptions
): Promise<unknown>;
