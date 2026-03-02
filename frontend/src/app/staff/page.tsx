import { listStaff } from "@/lib/api/staff";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Staff } from "@/types/staff";

const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID ?? "";

const ROLE_LABELS: Record<string, string> = {
  platform_admin: "Platform Admin",
  owner: "Owner",
  manager: "Manager",
  waiter: "Waiter",
  server: "Server",
  kitchen: "Kitchen Chef",
  bartender: "Bartender",
  cashier: "Cashier",
  warehouse: "Warehouse",
  accountant: "Accountant",
};

async function getStaff(): Promise<Staff[]> {
  if (!TENANT_ID) return [];
  try {
    return await listStaff(TENANT_ID);
  } catch {
    return [];
  }
}

function roleLabel(role: string): string {
  return ROLE_LABELS[role] ?? role;
}

export default async function StaffPage(): Promise<React.JSX.Element> {
  const staff = await getStaff();

  if (!TENANT_ID) {
    return (
      <div className="flex flex-col gap-4">
        <h2 className="text-2xl font-bold">Staff</h2>
        <p className="text-muted-foreground text-sm">
          Set <code className="font-mono bg-muted px-1 rounded">NEXT_PUBLIC_TENANT_ID</code> to view staff.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Staff</h2>
        <span className="text-sm text-muted-foreground">
          {staff.length} member{staff.length !== 1 ? "s" : ""}
        </span>
      </div>

      {staff.length === 0 ? (
        <p className="text-muted-foreground">No staff members found.</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {staff.map((member) => (
            <Card key={member.id}>
              <CardHeader>
                <div className="flex items-start justify-between gap-2">
                  <CardTitle className="text-base">{member.full_name}</CardTitle>
                  <Badge variant={member.is_active ? "success" : "outline"}>
                    {member.is_active ? "Active" : "Inactive"}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                <dl className="flex flex-col gap-1 text-sm">
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">Role</dt>
                    <dd className="font-medium capitalize">{roleLabel(member.role)}</dd>
                  </div>
                  {member.email && (
                    <div className="flex justify-between">
                      <dt className="text-muted-foreground">Email</dt>
                      <dd className="text-xs truncate max-w-[12rem] text-right">{member.email}</dd>
                    </div>
                  )}
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">Joined</dt>
                    <dd className="text-xs">
                      {new Date(member.created_at).toLocaleDateString("vi-VN")}
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
