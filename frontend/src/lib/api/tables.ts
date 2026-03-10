import { apiClient } from "./client";
import type { FloorPlan, RestaurantTable } from "@/types/table";

const TENANT_HEADER = (tenantID: string): HeadersInit => ({ "X-Tenant-ID": tenantID });

export interface CreateFloorPlanRequest {
  name: string;
  floor_level?: number;
}

export interface UpdateFloorPlanRequest {
  name?: string;
  floor_level?: number;
  is_active?: boolean;
}

export interface CreateTableRequest {
  floor_plan_id: string;
  name: string;
  capacity?: number;
  pos_x?: number;
  pos_y?: number;
  shape?: string;
}

export interface UpdateTableRequest {
  name?: string;
  capacity?: number;
  pos_x?: number;
  pos_y?: number;
  shape?: string;
  status?: string;
}

export function listFloorPlans(tenantID: string): Promise<FloorPlan[]> {
  return apiClient.get<FloorPlan[]>("/api/v1/floor-plans", { headers: TENANT_HEADER(tenantID) });
}

export function createFloorPlan(tenantID: string, data: CreateFloorPlanRequest): Promise<FloorPlan> {
  return apiClient.post<FloorPlan>("/api/v1/floor-plans", data, { headers: TENANT_HEADER(tenantID) });
}

export function updateFloorPlan(id: string, tenantID: string, data: UpdateFloorPlanRequest): Promise<FloorPlan> {
  return apiClient.put<FloorPlan>(`/api/v1/floor-plans/${id}`, data, { headers: TENANT_HEADER(tenantID) });
}

export function listTablesByFloorPlan(floorPlanID: string, tenantID: string): Promise<RestaurantTable[]> {
  return apiClient.get<RestaurantTable[]>(`/api/v1/floor-plans/${floorPlanID}/tables`, { headers: TENANT_HEADER(tenantID) });
}

export function createTable(tenantID: string, data: CreateTableRequest): Promise<RestaurantTable> {
  return apiClient.post<RestaurantTable>("/api/v1/tables", data, { headers: TENANT_HEADER(tenantID) });
}

export function updateTable(id: string, tenantID: string, data: UpdateTableRequest): Promise<RestaurantTable> {
  return apiClient.put<RestaurantTable>(`/api/v1/tables/${id}`, data, { headers: TENANT_HEADER(tenantID) });
}

export function deleteTable(id: string, tenantID: string): Promise<void> {
  return apiClient.delete<void>(`/api/v1/tables/${id}`, { headers: TENANT_HEADER(tenantID) });
}
