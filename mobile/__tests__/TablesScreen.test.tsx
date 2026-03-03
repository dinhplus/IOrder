import React from "react";
import { render, screen } from "@testing-library/react-native";
import TablesScreen from "../src/screens/TablesScreen";

const mockRefetch = jest.fn();

jest.mock("../src/hooks/useTables", () => ({
  useTables: jest.fn(),
}));

import { useTables } from "../src/hooks/useTables";

const mockUseTables = useTables as jest.MockedFunction<typeof useTables>;

describe("TablesScreen", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("shows loading indicator while fetching", () => {
    mockUseTables.mockReturnValue({
      data: [],
      loading: true,
      error: null,
      refetch: mockRefetch,
    });

    render(<TablesScreen />);
    expect(screen.getByTestId("activity-indicator")).toBeTruthy();
  });

  it("shows error message when request fails", () => {
    mockUseTables.mockReturnValue({
      data: [],
      loading: false,
      error: new Error("Network error"),
      refetch: mockRefetch,
    });

    render(<TablesScreen />);
    expect(screen.getByText("Unable to load floor plans")).toBeTruthy();
  });

  it("renders floor plans and tables when data is available", () => {
    mockUseTables.mockReturnValue({
      data: [
        {
          plan: { id: "plan-1", name: "Main Floor", floor_level: 1, is_active: true, tenant_id: "t", created_at: "" },
          tables: [
            { id: "table-1", name: "T1", capacity: 4, status: "available", floor_plan_id: "plan-1", tenant_id: "t", pos_x: 0, pos_y: 0, shape: "square", created_at: "", updated_at: "" },
          ],
        },
      ],
      loading: false,
      error: null,
      refetch: mockRefetch,
    });

    render(<TablesScreen />);
    expect(screen.getByText("Main Floor")).toBeTruthy();
    expect(screen.getByText("T1")).toBeTruthy();
  });

  it("renders empty state when no floor plans", () => {
    mockUseTables.mockReturnValue({
      data: [],
      loading: false,
      error: null,
      refetch: mockRefetch,
    });

    render(<TablesScreen />);
    expect(screen.getByText("No floor plans found.")).toBeTruthy();
  });
});
