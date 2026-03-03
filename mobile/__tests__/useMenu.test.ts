import { renderHook, act, waitFor } from "@testing-library/react-native";
import { useMenu } from "../src/hooks/useMenu";
import { APIError } from "../src/types/api";

const mockListCategories = jest.fn();
const mockListItems = jest.fn();

jest.mock("../src/lib/api/menu", () => ({
  listCategories: (...args: unknown[]) => mockListCategories(...args),
  listItems: (...args: unknown[]) => mockListItems(...args),
}));

jest.mock("../src/lib/config", () => ({
  getTenantId: () => "test-tenant",
}));

describe("useMenu", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("returns loading state initially", () => {
    mockListCategories.mockReturnValue(new Promise(() => {}));
    mockListItems.mockReturnValue(new Promise(() => {}));

    const { result } = renderHook(() => useMenu());
    expect(result.current.loading).toBe(true);
    expect(result.current.categories).toEqual([]);
    expect(result.current.items).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it("returns menu data on success", async () => {
    const mockCategories = [
      { id: "cat-1", name: "Starters", type: "food", is_active: true, sort_order: 0, tenant_id: "t", created_at: "" },
    ];
    const mockItems = [
      { id: "item-1", category_id: "cat-1", name: "Spring Rolls", price: 50000, is_available: true, sort_order: 0, tags: [], tenant_id: "t", created_at: "", updated_at: "" },
    ];
    mockListCategories.mockResolvedValue(mockCategories);
    mockListItems.mockResolvedValue(mockItems);

    const { result } = renderHook(() => useMenu());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.categories).toEqual(mockCategories);
    expect(result.current.items).toEqual(mockItems);
    expect(result.current.error).toBeNull();
  });

  it("returns error when request fails", async () => {
    const mockError = new APIError("INTERNAL", "Server error", undefined, 500);
    mockListCategories.mockRejectedValue(mockError);
    mockListItems.mockResolvedValue([]);

    const { result } = renderHook(() => useMenu());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.categories).toEqual([]);
    expect(result.current.error).toBe(mockError);
  });

  it("refetch re-triggers the request", async () => {
    mockListCategories.mockResolvedValue([]);
    mockListItems.mockResolvedValue([]);

    const { result } = renderHook(() => useMenu());
    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(mockListCategories).toHaveBeenCalledTimes(1);

    act(() => {
      result.current.refetch();
    });

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(mockListCategories).toHaveBeenCalledTimes(2);
  });
});
