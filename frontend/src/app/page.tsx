import { apiClient } from "@/lib/api/client";
import { type HealthResponse } from "@/types/health";
import { SystemStatus } from "@/components/features/system-status";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";

const FEATURES = [
  {
    title: "Customer Ordering",
    description: "Scan QR code to browse the menu and place orders from your table.",
    icon: "🍽️",
  },
  {
    title: "Kitchen Display",
    description: "Real-time order queue for kitchen and bar staff.",
    icon: "👨‍🍳",
  },
  {
    title: "Point of Sale",
    description: "Waiter POS with floor map and order management.",
    icon: "💳",
  },
  {
    title: "Owner Dashboard",
    description: "Analytics, menu management, and inventory tracking.",
    icon: "📊",
  },
  {
    title: "Membership",
    description: "Loyalty points, vouchers, and membership tiers.",
    icon: "🎁",
  },
  {
    title: "Payments",
    description: "Integrated payment processing and receipts.",
    icon: "🏦",
  },
];

async function getHealth(): Promise<HealthResponse | null> {
  try {
    return await apiClient.get<HealthResponse>("/health");
  } catch {
    return null;
  }
}

export default async function Home(): Promise<React.JSX.Element> {
  const health = await getHealth();

  return (
    <div className="flex flex-col gap-12">
      <section className="flex flex-col gap-4 py-8">
        <h2 className="text-4xl font-bold tracking-tight">Welcome to IOrder</h2>
        <p className="text-muted-foreground text-lg max-w-xl">
          A modern, multi-tenant restaurant ordering platform. Manage your restaurant from
          anywhere — orders, menu, payments, and more.
        </p>
      </section>

      <section className="flex flex-col gap-4">
        <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wider">
          System Status
        </h3>
        <SystemStatus health={health} />
      </section>

      <section className="flex flex-col gap-6">
        <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wider">
          Platform Features
        </h3>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {FEATURES.map((feature) => (
            <Card key={feature.title}>
              <CardHeader>
                <div className="text-3xl mb-2" aria-hidden="true">
                  {feature.icon}
                </div>
                <CardTitle>{feature.title}</CardTitle>
                <CardDescription>{feature.description}</CardDescription>
              </CardHeader>
              <CardContent>
                <span className="text-xs text-muted-foreground border border-border rounded px-2 py-0.5">
                  Coming soon
                </span>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>
    </div>
  );
}
