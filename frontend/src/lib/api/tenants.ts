import { apiClient } from "./client";
import type { Tenant } from "@/types/tenant";

export function listTenants(): Promise<Tenant[]> {
  return apiClient.get<Tenant[]>("/api/v1/tenants");
}

export function getTenant(id: string): Promise<Tenant> {
  return apiClient.get<Tenant>(`/api/v1/tenants/${id}`);
}
