import { apiClient } from "./client";
import type { Order } from "@/types/order";

const TENANT_HEADER = (tenantID: string): HeadersInit => ({ "X-Tenant-ID": tenantID });

export function listOrders(tenantID: string): Promise<Order[]> {
  return apiClient.get<Order[]>("/api/v1/orders", { headers: TENANT_HEADER(tenantID) });
}

export function getOrder(id: string): Promise<Order> {
  return apiClient.get<Order>(`/api/v1/orders/${id}`);
}

export function submitOrder(id: string): Promise<Order> {
  return apiClient.post<Order>(`/api/v1/orders/${id}/submit`, {});
}

export function confirmOrder(id: string): Promise<Order> {
  return apiClient.post<Order>(`/api/v1/orders/${id}/confirm`, {});
}

export function rejectOrder(id: string): Promise<Order> {
  return apiClient.post<Order>(`/api/v1/orders/${id}/reject`, {});
}

export function startPreparation(id: string): Promise<Order> {
  return apiClient.post<Order>(`/api/v1/orders/${id}/start-preparation`, {});
}

export function markReady(id: string): Promise<Order> {
  return apiClient.post<Order>(`/api/v1/orders/${id}/ready`, {});
}

export function markServed(id: string): Promise<Order> {
  return apiClient.post<Order>(`/api/v1/orders/${id}/serve`, {});
}

export function requestPayment(id: string): Promise<Order> {
  return apiClient.post<Order>(`/api/v1/orders/${id}/request-payment`, {});
}

export function cancelOrder(id: string): Promise<Order> {
  return apiClient.post<Order>(`/api/v1/orders/${id}/cancel`, {});
}
