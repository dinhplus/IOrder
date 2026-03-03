import { renderHook, act, waitFor } from "@testing-library/react-native";
import { useOrders } from "../src/hooks/useOrders";
import { APIError } from "../src/types/api";

const mockListOrders = jest.fn();

jest.mock("../src/lib/api/orders", () => ({
  listOrders: (...args: unknown[]) => mockListOrders(...args),
}));

jest.mock("../src/lib/config", () => ({
  getTenantId: () => "test-tenant",
}));

const mockOrder = {
  id: "order-1",
  tenant_id: "test-tenant",
  table_id: "table-1",
  session_id: "session-1",
  status: "SUBMITTED" as const,
  subtotal: 100000,
  discount_amount: 0,
  total: 100000,
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-01T00:00:00Z",
};

describe("useOrders", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("returns loading state initially", () => {
    mockListOrders.mockReturnValue(new Promise(() => {}));

    const { result } = renderHook(() => useOrders());
    expect(result.current.loading).toBe(true);
    expect(result.current.orders).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it("returns orders on success", async () => {
    mockListOrders.mockResolvedValue([mockOrder]);

    const { result } = renderHook(() => useOrders());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.orders).toEqual([mockOrder]);
    expect(result.current.error).toBeNull();
  });

  it("returns error when request fails", async () => {
    const mockError = new APIError("INTERNAL", "Server error", undefined, 500);
    mockListOrders.mockRejectedValue(mockError);

    const { result } = renderHook(() => useOrders());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.orders).toEqual([]);
    expect(result.current.error).toBe(mockError);
  });

  it("refetch re-triggers the request", async () => {
    mockListOrders.mockResolvedValue([mockOrder]);

    const { result } = renderHook(() => useOrders());
    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(mockListOrders).toHaveBeenCalledTimes(1);

    act(() => {
      result.current.refetch();
    });

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(mockListOrders).toHaveBeenCalledTimes(2);
  });
});
