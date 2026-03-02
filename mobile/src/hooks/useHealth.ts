import { useState, useEffect, useCallback } from "react";
import { apiClient } from "@/lib/api/client";
import { APIError } from "@/types/api";

export interface HealthStatus {
  status: string;
  version: string;
  db: string;
}

export interface UseHealthResult {
  health: HealthStatus | null;
  loading: boolean;
  error: APIError | Error | null;
  refetch: () => void;
}

export function useHealth(): UseHealthResult {
  const [health, setHealth] = useState<HealthStatus | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<APIError | Error | null>(null);

  const fetchHealth = useCallback((): void => {
    setLoading(true);
    setError(null);

    apiClient
      .get<HealthStatus>("/health")
      .then((data) => {
        setHealth(data);
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err : new Error(String(err)));
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    fetchHealth();
  }, [fetchHealth]);

  return { health, loading, error, refetch: fetchHealth };
}
