import { useState, useEffect, useCallback } from "react";
import { getTenantId } from "@/lib/config";
import { listFloorPlans, listTablesByFloorPlan } from "@/lib/api/tables";
import { APIError } from "@/types/api";
import type { FloorPlan, RestaurantTable } from "@/types/table";

export interface FloorPlanWithTables {
  plan: FloorPlan;
  tables: RestaurantTable[];
}

export interface UseTablesResult {
  data: FloorPlanWithTables[];
  loading: boolean;
  error: APIError | Error | null;
  refetch: () => void;
}

export function useTables(): UseTablesResult {
  const [data, setData] = useState<FloorPlanWithTables[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<APIError | Error | null>(null);

  const fetchTables = useCallback((): void => {
    const tenantId = getTenantId();
    if (!tenantId) {
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);

    listFloorPlans(tenantId)
      .then((plans) =>
        Promise.all(
          plans.map((plan) =>
            listTablesByFloorPlan(plan.id, tenantId).then((tables) => ({ plan, tables })),
          ),
        ),
      )
      .then(setData)
      .catch((err: unknown) => {
        setError(err instanceof Error ? err : new Error(String(err)));
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    fetchTables();
  }, [fetchTables]);

  return { data, loading, error, refetch: fetchTables };
}
