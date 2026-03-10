import { apiClient } from "./client";
import type { Tenant } from "@/types/tenant";

export interface CreateTenantRequest {
  slug: string;
  name: string;
  logo_url?: string;
  timezone?: string;
  currency?: string;
}

export interface UpdateTenantRequest {
  slug?: string;
  name?: string;
  logo_url?: string;
  timezone?: string;
  currency?: string;
  is_active?: boolean;
}

export function listTenants(): Promise<Tenant[]> {
  return apiClient.get<Tenant[]>("/api/v1/tenants");
}

export function getTenant(id: string): Promise<Tenant> {
  return apiClient.get<Tenant>(`/api/v1/tenants/${id}`);
}

export function createTenant(data: CreateTenantRequest): Promise<Tenant> {
  return apiClient.post<Tenant>("/api/v1/tenants", data);
}

export function updateTenant(id: string, data: UpdateTenantRequest): Promise<Tenant> {
  return apiClient.patch<Tenant>(`/api/v1/tenants/${id}`, data);
}
