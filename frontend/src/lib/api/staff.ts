import { apiClient } from "./client";
import type { Staff } from "@/types/staff";

const TENANT_HEADER = (tenantID: string): HeadersInit => ({ "X-Tenant-ID": tenantID });

export interface CreateStaffRequest {
  full_name: string;
  email?: string;
  role: string;
}

export interface UpdateStaffRequest {
  full_name?: string;
  email?: string;
  role?: string;
  is_active?: boolean;
}

export function listStaff(tenantID: string): Promise<Staff[]> {
  return apiClient.get<Staff[]>("/api/v1/staff", { headers: TENANT_HEADER(tenantID) });
}

export function getStaff(id: string, tenantID: string): Promise<Staff> {
  return apiClient.get<Staff>(`/api/v1/staff/${id}`, { headers: TENANT_HEADER(tenantID) });
}

export function createStaff(tenantID: string, data: CreateStaffRequest): Promise<Staff> {
  return apiClient.post<Staff>("/api/v1/staff", data, { headers: TENANT_HEADER(tenantID) });
}

export function updateStaff(id: string, tenantID: string, data: UpdateStaffRequest): Promise<Staff> {
  return apiClient.patch<Staff>(`/api/v1/staff/${id}`, data, { headers: TENANT_HEADER(tenantID) });
}

export function deleteStaff(id: string, tenantID: string): Promise<void> {
  return apiClient.delete<void>(`/api/v1/staff/${id}`, { headers: TENANT_HEADER(tenantID) });
}
