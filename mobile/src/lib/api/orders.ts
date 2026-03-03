import { apiClient, tenantHeader } from "./client";
import type { ApiResponse } from "@/types/api";
import type { Order } from "@/types/order";

export function listOrders(tenantID: string): Promise<Order[]> {
  return apiClient
    .get<ApiResponse<Order[]>>("/api/v1/orders", { headers: tenantHeader(tenantID) })
    .then((res) => res.data);
}
