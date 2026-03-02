export interface MenuCategory {
  id: string;
  tenant_id: string;
  name: string;
  type: string;
  description?: string | null;
  sort_order: number;
  is_active: boolean;
  created_at: string;
}

export interface MenuItem {
  id: string;
  tenant_id: string;
  category_id: string;
  name: string;
  description?: string | null;
  price: number;
  image_url?: string | null;
  tags: string[];
  is_available: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}
