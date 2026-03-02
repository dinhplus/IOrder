import { Badge } from "@/components/ui/badge";
import type { OrderStatus } from "@/types/order";

interface OrderStatusBadgeProps {
  status: OrderStatus;
}

function statusVariant(status: OrderStatus): "default" | "success" | "warning" | "danger" | "outline" {
  switch (status) {
    case "DRAFT": return "outline";
    case "SUBMITTED": return "warning";
    case "CONFIRMED": return "default";
    case "IN_PREPARATION": return "warning";
    case "READY": return "success";
    case "SERVED": return "success";
    case "PAYMENT_REQUESTED": return "warning";
    case "PAID": return "success";
    case "CANCELLED": return "danger";
    case "REJECTED": return "danger";
    default: return "outline";
  }
}

function statusLabel(status: OrderStatus): string {
  switch (status) {
    case "DRAFT": return "Draft";
    case "SUBMITTED": return "Submitted";
    case "CONFIRMED": return "Confirmed";
    case "IN_PREPARATION": return "In Preparation";
    case "READY": return "Ready";
    case "SERVED": return "Served";
    case "PAYMENT_REQUESTED": return "Payment Requested";
    case "PAID": return "Paid";
    case "CANCELLED": return "Cancelled";
    case "REJECTED": return "Rejected";
    default: return status;
  }
}

export function OrderStatusBadge({ status }: OrderStatusBadgeProps): React.JSX.Element {
  return <Badge variant={statusVariant(status)}>{statusLabel(status)}</Badge>;
}
