import { renderHook, act, waitFor } from "@testing-library/react-native";
import { useTables } from "../src/hooks/useTables";
import { APIError } from "../src/types/api";

const mockListFloorPlans = jest.fn();
const mockListTablesByFloorPlan = jest.fn();

jest.mock("../src/lib/api/tables", () => ({
  listFloorPlans: (...args: unknown[]) => mockListFloorPlans(...args),
  listTablesByFloorPlan: (...args: unknown[]) => mockListTablesByFloorPlan(...args),
}));

jest.mock("../src/lib/config", () => ({
  getTenantId: () => "test-tenant",
}));

const mockPlan = {
  id: "plan-1",
  tenant_id: "test-tenant",
  name: "Main Floor",
  floor_level: 1,
  is_active: true,
  created_at: "2024-01-01T00:00:00Z",
};

const mockTable = {
  id: "table-1",
  tenant_id: "test-tenant",
  floor_plan_id: "plan-1",
  name: "Table 1",
  capacity: 4,
  pos_x: 0,
  pos_y: 0,
  shape: "square",
  status: "available",
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-01T00:00:00Z",
};

describe("useTables", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("returns loading state initially", () => {
    mockListFloorPlans.mockReturnValue(new Promise(() => {}));

    const { result } = renderHook(() => useTables());
    expect(result.current.loading).toBe(true);
    expect(result.current.data).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it("returns floor plans with tables on success", async () => {
    mockListFloorPlans.mockResolvedValue([mockPlan]);
    mockListTablesByFloorPlan.mockResolvedValue([mockTable]);

    const { result } = renderHook(() => useTables());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.data).toEqual([{ plan: mockPlan, tables: [mockTable] }]);
    expect(result.current.error).toBeNull();
  });

  it("returns error when request fails", async () => {
    const mockError = new APIError("INTERNAL", "Server error", undefined, 500);
    mockListFloorPlans.mockRejectedValue(mockError);

    const { result } = renderHook(() => useTables());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.data).toEqual([]);
    expect(result.current.error).toBe(mockError);
  });

  it("refetch re-triggers the request", async () => {
    mockListFloorPlans.mockResolvedValue([mockPlan]);
    mockListTablesByFloorPlan.mockResolvedValue([mockTable]);

    const { result } = renderHook(() => useTables());
    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(mockListFloorPlans).toHaveBeenCalledTimes(1);

    act(() => {
      result.current.refetch();
    });

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(mockListFloorPlans).toHaveBeenCalledTimes(2);
  });
});
