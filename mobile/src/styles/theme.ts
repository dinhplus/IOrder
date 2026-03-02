export const colors = {
  primary: "#E63946",
  primaryDark: "#C1121F",
  background: "#FFFFFF",
  surface: "#F8F9FA",
  border: "#E9ECEF",
  text: "#212529",
  textSecondary: "#6C757D",
  textDisabled: "#ADB5BD",
  success: "#2A9D8F",
  warning: "#E9C46A",
  error: "#E63946",
  white: "#FFFFFF",
  black: "#000000",
} as const;

export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
} as const;

export const typography = {
  h1: { fontSize: 32, fontWeight: "bold" as const, lineHeight: 40 },
  h2: { fontSize: 24, fontWeight: "bold" as const, lineHeight: 32 },
  h3: { fontSize: 20, fontWeight: "600" as const, lineHeight: 28 },
  body: { fontSize: 16, fontWeight: "normal" as const, lineHeight: 24 },
  small: { fontSize: 14, fontWeight: "normal" as const, lineHeight: 20 },
  caption: { fontSize: 12, fontWeight: "normal" as const, lineHeight: 16 },
} as const;

export const borderRadius = {
  sm: 4,
  md: 8,
  lg: 12,
  xl: 16,
  full: 9999,
} as const;
