import { listTenants } from "@/lib/api/tenants";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import type { Tenant } from "@/types/tenant";

async function getTenants(): Promise<Tenant[]> {
  try {
    return await listTenants();
  } catch {
    return [];
  }
}

export default async function TenantsPage(): Promise<React.JSX.Element> {
  const tenants = await getTenants();

  return (
    <div className="flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Tenants</h2>
        <span className="text-sm text-muted-foreground">
          {tenants.length} restaurant{tenants.length !== 1 ? "s" : ""}
        </span>
      </div>

      {tenants.length === 0 ? (
        <p className="text-muted-foreground">No tenants found.</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {tenants.map((tenant) => (
            <Card key={tenant.id}>
              <CardHeader>
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <CardTitle className="text-base">{tenant.name}</CardTitle>
                    <CardDescription className="font-mono text-xs">{tenant.slug}</CardDescription>
                  </div>
                  <Badge variant={tenant.is_active ? "success" : "outline"}>
                    {tenant.is_active ? "Active" : "Inactive"}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                <dl className="flex flex-col gap-1 text-sm">
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">Currency</dt>
                    <dd className="font-medium">{tenant.currency}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">Timezone</dt>
                    <dd className="text-xs">{tenant.timezone}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">Created</dt>
                    <dd className="text-xs">
                      {new Date(tenant.created_at).toLocaleDateString("vi-VN")}
                    </dd>
                  </div>
                </dl>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
