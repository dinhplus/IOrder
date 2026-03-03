import { APIError, type APIErrorResponse } from "@/types/api";

const DEFAULT_API_URL = "http://localhost:8080";

function getApiUrl(): string {
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  const Constants = require("expo-constants").default as {
    expoConfig?: { extra?: { apiUrl?: string } };
  };
  return Constants.expoConfig?.extra?.apiUrl ?? DEFAULT_API_URL;
}

export function tenantHeader(tenantID: string): HeadersInit {
  return { "X-Tenant-ID": tenantID };
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const url = `${getApiUrl()}${path}`;
  const response = await fetch(url, {
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
    ...init,
  });

  if (!response.ok) {
    let errorBody: APIErrorResponse | undefined;
    try {
      errorBody = (await response.json()) as APIErrorResponse;
    } catch {
      throw new APIError("UNKNOWN_ERROR", response.statusText, undefined, response.status);
    }
    throw new APIError(
      errorBody?.error?.code ?? "UNKNOWN_ERROR",
      errorBody?.error?.message ?? response.statusText,
      errorBody?.error?.details,
      response.status,
    );
  }

  return response.json() as Promise<T>;
}

export const apiClient = {
  get: <T>(path: string, init?: Omit<RequestInit, "method">) =>
    request<T>(path, { ...init, method: "GET" }),
  post: <T>(path: string, body: unknown, init?: Omit<RequestInit, "method" | "body">) =>
    request<T>(path, { ...init, method: "POST", body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown, init?: Omit<RequestInit, "method" | "body">) =>
    request<T>(path, { ...init, method: "PUT", body: JSON.stringify(body) }),
  delete: <T>(path: string, init?: Omit<RequestInit, "method">) =>
    request<T>(path, { ...init, method: "DELETE" }),
};
