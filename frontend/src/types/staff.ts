export interface Staff {
  id: string;
  tenant_id: string;
  cognito_sub?: string | null;
  full_name: string;
  email?: string | null;
  role: string;
  is_active: boolean;
  created_at: string;
}
