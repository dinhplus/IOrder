// eslint-disable-next-line @typescript-eslint/no-var-requires
const Constants = require("expo-constants").default as {
  expoConfig?: { extra?: { tenantId?: string } };
};

export function getTenantId(): string {
  return Constants.expoConfig?.extra?.tenantId ?? "";
}
