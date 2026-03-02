export interface Tenant {
  id: string;
  slug: string;
  name: string;
  logo_url?: string | null;
  timezone: string;
  currency: string;
  is_active: boolean;
  settings: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}
