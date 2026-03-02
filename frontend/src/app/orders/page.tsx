import { listOrders } from "@/lib/api/orders";
import { OrderStatusBadge } from "@/components/features/orders/order-status-badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Order, OrderStatus } from "@/types/order";

const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID ?? "";

const STATUS_ORDER: OrderStatus[] = [
  "SUBMITTED",
  "CONFIRMED",
  "IN_PREPARATION",
  "READY",
  "PAYMENT_REQUESTED",
  "DRAFT",
  "SERVED",
  "PAID",
  "CANCELLED",
  "REJECTED",
];

async function getOrders(): Promise<Order[]> {
  if (!TENANT_ID) return [];
  try {
    return await listOrders(TENANT_ID);
  } catch {
    return [];
  }
}

export default async function OrdersPage(): Promise<React.JSX.Element> {
  const orders = await getOrders();

  if (!TENANT_ID) {
    return (
      <div className="flex flex-col gap-4">
        <h2 className="text-2xl font-bold">Orders</h2>
        <p className="text-muted-foreground text-sm">
          Set <code className="font-mono bg-muted px-1 rounded">NEXT_PUBLIC_TENANT_ID</code> to view orders.
        </p>
      </div>
    );
  }

  const grouped = STATUS_ORDER.reduce<Record<OrderStatus, Order[]>>((acc, s) => {
    acc[s] = orders.filter((o) => o.status === s);
    return acc;
  }, {} as Record<OrderStatus, Order[]>);

  const activeStatuses = STATUS_ORDER.filter((s) => grouped[s].length > 0);

  return (
    <div className="flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Orders</h2>
        <span className="text-sm text-muted-foreground">{orders.length} total</span>
      </div>

      {orders.length === 0 ? (
        <p className="text-muted-foreground">No orders found.</p>
      ) : (
        activeStatuses.map((status) => (
          <section key={status} className="flex flex-col gap-4">
            <div className="flex items-center gap-3">
              <OrderStatusBadge status={status} />
              <span className="text-sm text-muted-foreground">
                {grouped[status].length} order{grouped[status].length !== 1 ? "s" : ""}
              </span>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {grouped[status].map((order) => (
                <Card key={order.id}>
                  <CardHeader>
                    <CardTitle className="font-mono text-sm">
                      #{order.id.slice(0, 8)}
                    </CardTitle>
                    <p className="text-xs text-muted-foreground">
                      Table: {order.table_id.slice(0, 8)}
                    </p>
                  </CardHeader>
                  <CardContent>
                    <dl className="flex flex-col gap-1 text-sm">
                      <div className="flex justify-between">
                        <dt className="text-muted-foreground">Total</dt>
                        <dd className="font-medium">{order.total.toLocaleString("vi-VN")}₫</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-muted-foreground">Created</dt>
                        <dd className="text-xs">
                          {new Date(order.created_at).toLocaleString("vi-VN")}
                        </dd>
                      </div>
                      {order.notes && (
                        <div className="flex justify-between">
                          <dt className="text-muted-foreground">Notes</dt>
                          <dd className="text-xs max-w-[12rem] text-right truncate">{order.notes}</dd>
                        </div>
                      )}
                    </dl>
                  </CardContent>
                </Card>
              ))}
            </div>
          </section>
        ))
      )}
    </div>
  );
}
