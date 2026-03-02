import { type HealthResponse, type HealthStatus } from "@/types/health";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface SystemStatusProps {
  health: HealthResponse | null;
}

function statusVariant(
  status: HealthStatus | undefined,
): "success" | "warning" | "danger" | "outline" {
  if (!status) return "outline";
  if (status === "ok") return "success";
  if (status === "degraded") return "warning";
  return "danger";
}

function statusLabel(status: HealthStatus | undefined): string {
  if (!status) return "Unknown";
  if (status === "ok") return "Healthy";
  if (status === "degraded") return "Degraded";
  return "Unavailable";
}

export function SystemStatus({ health }: SystemStatusProps): React.JSX.Element {
  const isOffline = health === null;

  return (
    <Card className="max-w-sm">
      <CardHeader>
        <CardTitle>System Status</CardTitle>
      </CardHeader>
      <CardContent>
        <dl className="flex flex-col gap-3 text-sm">
          <div className="flex items-center justify-between">
            <dt className="text-muted-foreground">API</dt>
            <dd>
              <Badge variant={isOffline ? "danger" : statusVariant(health?.status)}>
                {isOffline ? "Offline" : statusLabel(health?.status)}
              </Badge>
            </dd>
          </div>
          <div className="flex items-center justify-between">
            <dt className="text-muted-foreground">Database</dt>
            <dd>
              <Badge variant={isOffline ? "outline" : statusVariant(health?.db)}>
                {isOffline ? "—" : statusLabel(health?.db)}
              </Badge>
            </dd>
          </div>
          {health?.version && (
            <div className="flex items-center justify-between">
              <dt className="text-muted-foreground">Version</dt>
              <dd className="font-mono text-xs">{health.version}</dd>
            </div>
          )}
        </dl>
      </CardContent>
    </Card>
  );
}
