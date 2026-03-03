import { View, Text, StyleSheet } from "react-native";
import { colors, borderRadius, typography } from "@/styles/theme";

export type BadgeVariant = "success" | "error" | "warning" | "info" | "default";

interface BadgeProps {
  label: string;
  variant?: BadgeVariant;
}

const variantStyles: Record<BadgeVariant, { bg: string; text: string }> = {
  success: { bg: colors.success + "22", text: colors.success },
  error: { bg: colors.error + "22", text: colors.error },
  warning: { bg: colors.warning + "44", text: colors.warningText },
  info: { bg: colors.info + "22", text: colors.infoText },
  default: { bg: colors.border, text: colors.textSecondary },
};

export default function Badge({
  label,
  variant = "default",
}: BadgeProps): React.JSX.Element {
  const vs = variantStyles[variant];
  return (
    <View style={[styles.badge, { backgroundColor: vs.bg }]}>
      <Text style={[styles.label, { color: vs.text }]}>{label}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  badge: {
    borderRadius: borderRadius.full,
    paddingHorizontal: 8,
    paddingVertical: 2,
    alignSelf: "flex-start",
  },
  label: {
    ...typography.caption,
    fontWeight: "600",
  },
});
