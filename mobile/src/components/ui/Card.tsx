import { View, StyleSheet, type ViewProps } from "react-native";
import { colors, spacing, borderRadius } from "@/styles/theme";

interface CardProps extends ViewProps {
  children: React.ReactNode;
  padded?: boolean;
}

export default function Card({
  children,
  padded = true,
  style,
  ...rest
}: CardProps): React.JSX.Element {
  return (
    <View style={[styles.card, padded && styles.padded, style]} {...rest}>
      {children}
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.background,
    borderRadius: borderRadius.lg,
    borderWidth: 1,
    borderColor: colors.border,
    shadowColor: colors.black,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 4,
    elevation: 2,
  },
  padded: {
    padding: spacing.md,
  },
});
