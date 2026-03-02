import { renderHook, act, waitFor } from "@testing-library/react-native";
import { useHealth } from "../src/hooks/useHealth";
import { APIError } from "../src/types/api";

const mockGet = jest.fn();

jest.mock("../src/lib/api/client", () => ({
  apiClient: {
    get: (...args: unknown[]) => mockGet(...args),
  },
}));

describe("useHealth", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("returns loading state initially", () => {
    mockGet.mockReturnValue(new Promise(() => {}));
    const { result } = renderHook(() => useHealth());
    expect(result.current.loading).toBe(true);
    expect(result.current.health).toBeNull();
    expect(result.current.error).toBeNull();
  });

  it("returns health data on success", async () => {
    const mockData = { status: "ok", version: "0.1.0", db: "ok" };
    mockGet.mockResolvedValue(mockData);

    const { result } = renderHook(() => useHealth());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.health).toEqual(mockData);
    expect(result.current.error).toBeNull();
  });

  it("returns error when request fails", async () => {
    const mockError = new APIError("INTERNAL", "Server error", undefined, 500);
    mockGet.mockRejectedValue(mockError);

    const { result } = renderHook(() => useHealth());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.health).toBeNull();
    expect(result.current.error).toBe(mockError);
  });

  it("refetch re-triggers the request", async () => {
    const mockData = { status: "ok", version: "0.1.0", db: "ok" };
    mockGet.mockResolvedValue(mockData);

    const { result } = renderHook(() => useHealth());
    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(mockGet).toHaveBeenCalledTimes(1);

    act(() => {
      result.current.refetch();
    });

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(mockGet).toHaveBeenCalledTimes(2);
  });
});
