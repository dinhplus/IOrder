import React from "react";
import { render, screen } from "@testing-library/react-native";
import OrdersScreen from "../src/screens/OrdersScreen";

const mockRefetch = jest.fn();

jest.mock("../src/hooks/useOrders", () => ({
  useOrders: jest.fn(),
}));

import { useOrders } from "../src/hooks/useOrders";

const mockUseOrders = useOrders as jest.MockedFunction<typeof useOrders>;

const mockOrder = {
  id: "abcdef12-0000-0000-0000-000000000000",
  tenant_id: "test-tenant",
  table_id: "table-1-0000-0000-0000-000000000000",
  session_id: "session-1",
  status: "SUBMITTED" as const,
  subtotal: 100000,
  discount_amount: 0,
  total: 100000,
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-01T00:00:00Z",
};

describe("OrdersScreen", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("shows loading indicator while fetching", () => {
    mockUseOrders.mockReturnValue({
      orders: [],
      loading: true,
      error: null,
      refetch: mockRefetch,
    });

    render(<OrdersScreen />);
    expect(screen.getByTestId("activity-indicator")).toBeTruthy();
  });

  it("shows error message when request fails", () => {
    mockUseOrders.mockReturnValue({
      orders: [],
      loading: false,
      error: new Error("Network error"),
      refetch: mockRefetch,
    });

    render(<OrdersScreen />);
    expect(screen.getByText("Unable to load orders")).toBeTruthy();
  });

  it("renders orders when data is available", () => {
    mockUseOrders.mockReturnValue({
      orders: [mockOrder],
      loading: false,
      error: null,
      refetch: mockRefetch,
    });

    render(<OrdersScreen />);
    expect(screen.getByText("#abcdef12")).toBeTruthy();
  });

  it("renders empty state when no orders", () => {
    mockUseOrders.mockReturnValue({
      orders: [],
      loading: false,
      error: null,
      refetch: mockRefetch,
    });

    render(<OrdersScreen />);
    expect(screen.getByText("No orders found.")).toBeTruthy();
  });
});
