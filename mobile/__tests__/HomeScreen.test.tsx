import React from "react";
import { render, screen } from "@testing-library/react-native";
import HomeScreen from "../src/screens/HomeScreen";

const mockRefetch = jest.fn();

jest.mock("../src/hooks/useHealth", () => ({
  useHealth: jest.fn(),
}));

import { useHealth } from "../src/hooks/useHealth";

const mockUseHealth = useHealth as jest.MockedFunction<typeof useHealth>;

describe("HomeScreen", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("shows loading indicator while fetching", () => {
    mockUseHealth.mockReturnValue({
      health: null,
      loading: true,
      error: null,
      refetch: mockRefetch,
    });

    render(<HomeScreen />);
    expect(screen.getByText("Checking connection…")).toBeTruthy();
  });

  it("shows error message when request fails", () => {
    mockUseHealth.mockReturnValue({
      health: null,
      loading: false,
      error: new Error("Network error"),
      refetch: mockRefetch,
    });

    render(<HomeScreen />);
    expect(screen.getByText("Unable to reach the server")).toBeTruthy();
  });

  it("renders health status when data is available", () => {
    mockUseHealth.mockReturnValue({
      health: { status: "ok", version: "0.1.0", db: "ok" },
      loading: false,
      error: null,
      refetch: mockRefetch,
    });

    render(<HomeScreen />);
    expect(screen.getAllByText("ok")).toHaveLength(2);
    expect(screen.getByText("0.1.0")).toBeTruthy();
  });

  it("renders title and subtitle", () => {
    mockUseHealth.mockReturnValue({
      health: null,
      loading: true,
      error: null,
      refetch: mockRefetch,
    });

    render(<HomeScreen />);
    expect(screen.getByText("IOrder")).toBeTruthy();
    expect(screen.getByText("A modern restaurant ordering platform")).toBeTruthy();
  });
});
