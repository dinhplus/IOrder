import { apiClient, tenantHeader } from "./client";
import type { ApiResponse } from "@/types/api";
import type { MenuCategory, MenuItem } from "@/types/menu";

export function listCategories(tenantID: string): Promise<MenuCategory[]> {
  return apiClient
    .get<ApiResponse<MenuCategory[]>>("/api/v1/menu/categories", { headers: tenantHeader(tenantID) })
    .then((res) => res.data);
}

export function listItems(tenantID: string): Promise<MenuItem[]> {
  return apiClient
    .get<ApiResponse<MenuItem[]>>("/api/v1/menu/items", { headers: tenantHeader(tenantID) })
    .then((res) => res.data);
}
