export type HealthStatus = "ok" | "degraded" | "unavailable";

export interface HealthResponse {
  status: HealthStatus;
  version: string;
  db?: HealthStatus;
}
