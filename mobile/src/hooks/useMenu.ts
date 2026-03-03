import { useState, useEffect, useCallback } from "react";
import { getTenantId } from "@/lib/config";
import { listCategories, listItems } from "@/lib/api/menu";
import { APIError } from "@/types/api";
import type { MenuCategory, MenuItem } from "@/types/menu";

export interface UseMenuResult {
  categories: MenuCategory[];
  items: MenuItem[];
  loading: boolean;
  error: APIError | Error | null;
  refetch: () => void;
}

export function useMenu(): UseMenuResult {
  const [categories, setCategories] = useState<MenuCategory[]>([]);
  const [items, setItems] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<APIError | Error | null>(null);

  const fetchMenu = useCallback((): void => {
    const tenantId = getTenantId();
    if (!tenantId) {
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);

    Promise.all([listCategories(tenantId), listItems(tenantId)])
      .then(([cats, itms]) => {
        setCategories(cats);
        setItems(itms);
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err : new Error(String(err)));
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    fetchMenu();
  }, [fetchMenu]);

  return { categories, items, loading, error, refetch: fetchMenu };
}
