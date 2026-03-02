import { apiClient } from "./client";
import type { FloorPlan, RestaurantTable } from "@/types/table";

const TENANT_HEADER = (tenantID: string): HeadersInit => ({ "X-Tenant-ID": tenantID });

export function listFloorPlans(tenantID: string): Promise<FloorPlan[]> {
  return apiClient.get<FloorPlan[]>("/api/v1/floor-plans", { headers: TENANT_HEADER(tenantID) });
}

export function listTablesByFloorPlan(floorPlanID: string, tenantID: string): Promise<RestaurantTable[]> {
  return apiClient.get<RestaurantTable[]>(`/api/v1/floor-plans/${floorPlanID}/tables`, { headers: TENANT_HEADER(tenantID) });
}
