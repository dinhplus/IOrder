import React from "react";
import { render, screen } from "@testing-library/react-native";
import MenuScreen from "../src/screens/MenuScreen";

const mockRefetch = jest.fn();

jest.mock("../src/hooks/useMenu", () => ({
  useMenu: jest.fn(),
}));

import { useMenu } from "../src/hooks/useMenu";

const mockUseMenu = useMenu as jest.MockedFunction<typeof useMenu>;

describe("MenuScreen", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("shows loading indicator while fetching", () => {
    mockUseMenu.mockReturnValue({
      categories: [],
      items: [],
      loading: true,
      error: null,
      refetch: mockRefetch,
    });

    render(<MenuScreen />);
    expect(screen.getByTestId("activity-indicator")).toBeTruthy();
  });

  it("shows error message when request fails", () => {
    mockUseMenu.mockReturnValue({
      categories: [],
      items: [],
      loading: false,
      error: new Error("Network error"),
      refetch: mockRefetch,
    });

    render(<MenuScreen />);
    expect(screen.getByText("Unable to load menu")).toBeTruthy();
  });

  it("renders categories and items when data is available", () => {
    mockUseMenu.mockReturnValue({
      categories: [
        { id: "cat-1", name: "Starters", type: "food", is_active: true, sort_order: 0, tenant_id: "t", created_at: "" },
      ],
      items: [
        { id: "item-1", category_id: "cat-1", name: "Spring Rolls", price: 50000, is_available: true, sort_order: 0, tags: [], tenant_id: "t", created_at: "", updated_at: "" },
      ],
      loading: false,
      error: null,
      refetch: mockRefetch,
    });

    render(<MenuScreen />);
    expect(screen.getByText("Starters")).toBeTruthy();
    expect(screen.getByText("Spring Rolls")).toBeTruthy();
  });

  it("renders empty state when no categories", () => {
    mockUseMenu.mockReturnValue({
      categories: [],
      items: [],
      loading: false,
      error: null,
      refetch: mockRefetch,
    });

    render(<MenuScreen />);
    expect(screen.getByText("No menu categories found.")).toBeTruthy();
  });
});
