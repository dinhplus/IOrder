export type OrderStatus =
  | "DRAFT"
  | "SUBMITTED"
  | "CONFIRMED"
  | "IN_PREPARATION"
  | "READY"
  | "SERVED"
  | "PAYMENT_REQUESTED"
  | "PAID"
  | "CLOSED"
  | "CANCELLED"
  | "REJECTED";

export const ORDER_STATUS_DISPLAY: OrderStatus[] = [
  "SUBMITTED",
  "CONFIRMED",
  "IN_PREPARATION",
  "READY",
  "PAYMENT_REQUESTED",
  "DRAFT",
  "SERVED",
  "PAID",
  "CLOSED",
  "CANCELLED",
  "REJECTED",
];

export interface OrderItem {
  id: string;
  order_id: string;
  tenant_id: string;
  item_id: string;
  item_name: string;
  item_price: number;
  quantity: number;
  modifiers: unknown;
  notes?: string | null;
  status: string;
  routed_to: string;
}

export interface Order {
  id: string;
  tenant_id: string;
  table_id: string;
  session_id: string;
  status: OrderStatus;
  customer_id?: string | null;
  subtotal: number;
  discount_amount: number;
  total: number;
  notes?: string | null;
  placed_at?: string | null;
  confirmed_at?: string | null;
  ready_at?: string | null;
  served_at?: string | null;
  paid_at?: string | null;
  created_at: string;
  updated_at: string;
}
