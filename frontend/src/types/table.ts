export interface FloorPlan {
  id: string;
  tenant_id: string;
  name: string;
  floor_level: number;
  is_active: boolean;
  created_at: string;
}

export interface RestaurantTable {
  id: string;
  tenant_id: string;
  floor_plan_id: string;
  name: string;
  capacity: number;
  pos_x: number;
  pos_y: number;
  shape: string;
  status: string;
  created_at: string;
  updated_at: string;
}
