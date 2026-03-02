import { listFloorPlans, listTablesByFloorPlan } from "@/lib/api/tables";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import type { FloorPlan, RestaurantTable } from "@/types/table";

const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID ?? "";

interface FloorPlanWithTables {
  plan: FloorPlan;
  tables: RestaurantTable[];
}

async function getData(): Promise<FloorPlanWithTables[]> {
  if (!TENANT_ID) return [];
  try {
    const plans = await listFloorPlans(TENANT_ID);
    const results = await Promise.all(
      plans.map(async (plan) => {
        const tables = await listTablesByFloorPlan(plan.id, TENANT_ID);
        return { plan, tables };
      }),
    );
    return results;
  } catch {
    return [];
  }
}

function tableStatusVariant(status: string): "success" | "warning" | "danger" | "outline" {
  if (status === "available") return "success";
  if (status === "occupied") return "danger";
  if (status === "reserved") return "warning";
  return "outline";
}

export default async function FloorPlanPage(): Promise<React.JSX.Element> {
  const data = await getData();

  if (!TENANT_ID) {
    return (
      <div className="flex flex-col gap-4">
        <h2 className="text-2xl font-bold">Floor Plan</h2>
        <p className="text-muted-foreground text-sm">
          Set <code className="font-mono bg-muted px-1 rounded">NEXT_PUBLIC_TENANT_ID</code> to view floor plans.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Floor Plan</h2>
        <span className="text-sm text-muted-foreground">
          {data.length} floor plan{data.length !== 1 ? "s" : ""}
        </span>
      </div>

      {data.length === 0 ? (
        <p className="text-muted-foreground">No floor plans found.</p>
      ) : (
        data.map(({ plan, tables }) => (
          <section key={plan.id} className="flex flex-col gap-4">
            <div className="flex items-center gap-3">
              <h3 className="text-lg font-semibold">{plan.name}</h3>
              <span className="text-xs text-muted-foreground">Floor {plan.floor_level}</span>
              <Badge variant={plan.is_active ? "success" : "outline"}>
                {plan.is_active ? "Active" : "Inactive"}
              </Badge>
              <span className="text-xs text-muted-foreground">{tables.length} tables</span>
            </div>
            {tables.length === 0 ? (
              <p className="text-sm text-muted-foreground pl-4">No tables on this floor plan.</p>
            ) : (
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
                {tables.map((table) => (
                  <Card key={table.id} className="text-center">
                    <CardHeader className="pb-1 px-3 pt-4">
                      <CardTitle className="text-sm">{table.name}</CardTitle>
                      <CardDescription className="text-xs">{table.capacity} seats</CardDescription>
                    </CardHeader>
                    <CardContent className="pb-4 px-3 pt-0">
                      <Badge variant={tableStatusVariant(table.status)} className="capitalize">
                        {table.status}
                      </Badge>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </section>
        ))
      )}
    </div>
  );
}
