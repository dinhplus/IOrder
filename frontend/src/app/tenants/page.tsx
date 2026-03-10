import { TenantsView } from "@/components/features/tenants/tenants-view";
import { listTenants } from "@/lib/api/tenants";

async function getTenants() {
  try {
    return await listTenants();
  } catch {
    return [];
  }
}

export default async function TenantsPage() {
  const initialTenants = await getTenants();

  return (
    <div className="flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Tenants</h2>
        <span className="text-sm text-muted-foreground">
          {initialTenants.length} restaurant{initialTenants.length !== 1 ? "s" : ""}
        </span>
      </div>
      <TenantsView initialTenants={initialTenants} />
    </div>
  );
}
