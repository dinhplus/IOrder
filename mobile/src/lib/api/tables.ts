import { apiClient, tenantHeader } from "./client";
import type { ApiResponse } from "@/types/api";
import type { FloorPlan, RestaurantTable } from "@/types/table";

export function listFloorPlans(tenantID: string): Promise<FloorPlan[]> {
  return apiClient
    .get<ApiResponse<FloorPlan[]>>("/api/v1/floor-plans", { headers: tenantHeader(tenantID) })
    .then((res) => res.data);
}

export function listTablesByFloorPlan(
  floorPlanID: string,
  tenantID: string,
): Promise<RestaurantTable[]> {
  return apiClient
    .get<ApiResponse<RestaurantTable[]>>(`/api/v1/floor-plans/${floorPlanID}/tables`, {
      headers: tenantHeader(tenantID),
    })
    .then((res) => res.data);
}
