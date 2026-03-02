import { apiClient } from "./client";
import type { Staff } from "@/types/staff";

const TENANT_HEADER = (tenantID: string): HeadersInit => ({ "X-Tenant-ID": tenantID });

export function listStaff(tenantID: string): Promise<Staff[]> {
  return apiClient.get<Staff[]>("/api/v1/staff", { headers: TENANT_HEADER(tenantID) });
}

export function getStaff(id: string, tenantID: string): Promise<Staff> {
  return apiClient.get<Staff>(`/api/v1/staff/${id}`, { headers: TENANT_HEADER(tenantID) });
}
