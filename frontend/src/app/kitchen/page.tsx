import { listOrders } from "@/lib/api/orders";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { Order } from "@/types/order";

const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID ?? "";

async function getOrders(): Promise<Order[]> {
  if (!TENANT_ID) return [];
  try {
    return await listOrders(TENANT_ID);
  } catch {
    return [];
  }
}

interface KDSLaneProps {
  title: string;
  orders: Order[];
  badgeVariant: "default" | "success" | "warning" | "danger" | "outline";
  emptyText: string;
}

function KDSLane({ title, orders, badgeVariant, emptyText }: KDSLaneProps): React.JSX.Element {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-3 pb-2 border-b border-border">
        <h3 className="text-lg font-semibold">{title}</h3>
        <Badge variant={badgeVariant}>{orders.length}</Badge>
      </div>
      {orders.length === 0 ? (
        <p className="text-sm text-muted-foreground">{emptyText}</p>
      ) : (
        <div className="flex flex-col gap-3">
          {orders.map((order) => (
            <Card key={order.id}>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-mono">#{order.id.slice(0, 8)}</CardTitle>
                <p className="text-xs text-muted-foreground">
                  Table: {order.table_id.slice(0, 8)}
                </p>
              </CardHeader>
              <CardContent className="pt-0">
                {order.notes && (
                  <p className="text-xs text-muted-foreground italic">&quot;{order.notes}&quot;</p>
                )}
                <p className="text-xs text-muted-foreground mt-1">
                  {new Date(order.created_at).toLocaleTimeString("vi-VN")}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

export default async function KitchenPage(): Promise<React.JSX.Element> {
  const orders = await getOrders();

  const newOrders = orders.filter((o) => o.status === "CONFIRMED");
  const inProgressOrders = orders.filter((o) => o.status === "IN_PREPARATION");
  const readyOrders = orders.filter((o) => o.status === "READY");

  return (
    <div className="flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Kitchen Display</h2>
        {!TENANT_ID && (
          <p className="text-sm text-muted-foreground">
            Set <code className="font-mono bg-muted px-1 rounded">NEXT_PUBLIC_TENANT_ID</code> to view orders.
          </p>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <KDSLane
          title="New"
          orders={newOrders}
          badgeVariant="warning"
          emptyText="No new orders"
        />
        <KDSLane
          title="In Progress"
          orders={inProgressOrders}
          badgeVariant="default"
          emptyText="Nothing in preparation"
        />
        <KDSLane
          title="Ready"
          orders={readyOrders}
          badgeVariant="success"
          emptyText="Nothing ready yet"
        />
      </div>
    </div>
  );
}
