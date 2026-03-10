import { apiClient } from "./client";
import type { MenuCategory, MenuItem } from "@/types/menu";

const TENANT_HEADER = (tenantID: string): HeadersInit => ({ "X-Tenant-ID": tenantID });

export interface CreateCategoryRequest {
  name: string;
  type?: string;
  description?: string;
  sort_order?: number;
}

export interface UpdateCategoryRequest {
  name?: string;
  type?: string;
  description?: string;
  sort_order?: number;
  is_active?: boolean;
}

export interface CreateMenuItemRequest {
  category_id: string;
  name: string;
  description?: string;
  price: number;
  image_url?: string;
  sort_order?: number;
}

export interface UpdateMenuItemRequest {
  category_id?: string;
  name?: string;
  description?: string;
  price?: number;
  image_url?: string;
  is_available?: boolean;
  sort_order?: number;
}

export function listCategories(tenantID: string): Promise<MenuCategory[]> {
  return apiClient.get<MenuCategory[]>("/api/v1/menu/categories", { headers: TENANT_HEADER(tenantID) });
}

export function createCategory(tenantID: string, data: CreateCategoryRequest): Promise<MenuCategory> {
  return apiClient.post<MenuCategory>("/api/v1/menu/categories", data, { headers: TENANT_HEADER(tenantID) });
}

export function updateCategory(id: string, tenantID: string, data: UpdateCategoryRequest): Promise<MenuCategory> {
  return apiClient.put<MenuCategory>(`/api/v1/menu/categories/${id}`, data, { headers: TENANT_HEADER(tenantID) });
}

export function deleteCategory(id: string, tenantID: string): Promise<void> {
  return apiClient.delete<void>(`/api/v1/menu/categories/${id}`, { headers: TENANT_HEADER(tenantID) });
}

export function listItems(tenantID: string): Promise<MenuItem[]> {
  return apiClient.get<MenuItem[]>("/api/v1/menu/items", { headers: TENANT_HEADER(tenantID) });
}

export function createMenuItem(tenantID: string, data: CreateMenuItemRequest): Promise<MenuItem> {
  return apiClient.post<MenuItem>("/api/v1/menu/items", data, { headers: TENANT_HEADER(tenantID) });
}

export function updateMenuItem(id: string, tenantID: string, data: UpdateMenuItemRequest): Promise<MenuItem> {
  return apiClient.put<MenuItem>(`/api/v1/menu/items/${id}`, data, { headers: TENANT_HEADER(tenantID) });
}

export function deleteMenuItem(id: string, tenantID: string): Promise<void> {
  return apiClient.delete<void>(`/api/v1/menu/items/${id}`, { headers: TENANT_HEADER(tenantID) });
}

export function toggleItemAvailability(id: string, tenantID: string, isAvailable: boolean): Promise<{ id: string; is_available: boolean }> {
  return apiClient.patch<{ id: string; is_available: boolean }>(
    `/api/v1/menu/items/${id}/availability`,
    { is_available: isAvailable },
    { headers: TENANT_HEADER(tenantID) },
  );
}
