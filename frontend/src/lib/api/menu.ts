import { apiClient } from "./client";
import type { MenuCategory, MenuItem } from "@/types/menu";

const TENANT_HEADER = (tenantID: string): HeadersInit => ({ "X-Tenant-ID": tenantID });

export function listCategories(tenantID: string): Promise<MenuCategory[]> {
  return apiClient.get<MenuCategory[]>("/api/v1/menu/categories", { headers: TENANT_HEADER(tenantID) });
}

export function listItems(tenantID: string): Promise<MenuItem[]> {
  return apiClient.get<MenuItem[]>("/api/v1/menu/items", { headers: TENANT_HEADER(tenantID) });
}

export function toggleItemAvailability(id: string, tenantID: string, isAvailable: boolean): Promise<{ id: string; is_available: boolean }> {
  return apiClient.patch<{ id: string; is_available: boolean }>(
    `/api/v1/menu/items/${id}/availability`,
    { is_available: isAvailable },
    { headers: TENANT_HEADER(tenantID) },
  );
}
