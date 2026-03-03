import { useState, useEffect, useCallback } from "react";
import { getTenantId } from "@/lib/config";
import { listOrders } from "@/lib/api/orders";
import { APIError } from "@/types/api";
import type { Order } from "@/types/order";

export interface UseOrdersResult {
  orders: Order[];
  loading: boolean;
  error: APIError | Error | null;
  refetch: () => void;
}

export function useOrders(): UseOrdersResult {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<APIError | Error | null>(null);

  const fetchOrders = useCallback((): void => {
    const tenantId = getTenantId();
    if (!tenantId) {
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);

    listOrders(tenantId)
      .then((data) => {
        setOrders(data);
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err : new Error(String(err)));
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    fetchOrders();
  }, [fetchOrders]);

  return { orders, loading, error, refetch: fetchOrders };
}
