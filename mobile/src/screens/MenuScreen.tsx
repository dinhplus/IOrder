import { View, Text, StyleSheet, ActivityIndicator, ScrollView } from "react-native";
import { StatusBar } from "expo-status-bar";
import { useMenu } from "@/hooks/useMenu";
import Card from "@/components/ui/Card";
import Badge from "@/components/ui/Badge";
import Button from "@/components/ui/Button";
import { formatCurrency } from "@/lib/format";
import { colors, spacing, typography } from "@/styles/theme";
import type { MenuItem } from "@/types/menu";

export default function MenuScreen(): React.JSX.Element {
  const { categories, items, loading, error, refetch } = useMenu();

  if (loading) {
    return (
      <View style={styles.centered}>
        <ActivityIndicator testID="activity-indicator" size="large" color={colors.primary} />
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.centered}>
        <Text style={styles.errorText}>Unable to load menu</Text>
        <Button title="Retry" onPress={refetch} style={styles.retryButton} />
      </View>
    );
  }

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <StatusBar style="auto" />

      <View style={styles.header}>
        <Text style={styles.title}>Menu</Text>
        <Text style={styles.subtitle}>
          {items.length} items · {categories.length} categories
        </Text>
      </View>

      {categories.length === 0 ? (
        <Text style={styles.emptyText}>No menu categories found.</Text>
      ) : (
        categories.map((cat) => {
          const catItems = items.filter((item) => item.category_id === cat.id);
          return (
            <View key={cat.id} style={styles.section}>
              <View style={styles.sectionHeader}>
                <Text style={styles.sectionTitle}>{cat.name}</Text>
                <Badge
                  label={cat.is_active ? "Active" : "Inactive"}
                  variant={cat.is_active ? "success" : "default"}
                />
                <Text style={styles.typeLabel}>{cat.type}</Text>
              </View>

              {catItems.length === 0 ? (
                <Text style={styles.emptyText}>No items in this category.</Text>
              ) : (
                catItems.map((item: MenuItem) => (
                  <Card key={item.id} style={styles.itemCard}>
                    <View style={styles.itemHeader}>
                      <Text style={styles.itemName}>{item.name}</Text>
                      <Badge
                        label={item.is_available ? "Available" : "Unavailable"}
                        variant={item.is_available ? "success" : "error"}
                      />
                    </View>
                    {item.description ? (
                      <Text style={styles.itemDescription}>{item.description}</Text>
                    ) : null}
                    <Text style={styles.itemPrice}>{formatCurrency(item.price)}</Text>
                  </Card>
                ))
              )}
            </View>
          );
        })
      )}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.surface,
  },
  content: {
    padding: spacing.lg,
    paddingBottom: spacing.xxl,
  },
  centered: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.surface,
    padding: spacing.lg,
  },
  header: {
    marginBottom: spacing.xl,
  },
  title: {
    ...typography.h2,
    color: colors.text,
  },
  subtitle: {
    ...typography.small,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  section: {
    marginBottom: spacing.xl,
  },
  sectionHeader: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: spacing.md,
    flexWrap: "wrap",
    gap: spacing.sm,
  },
  sectionTitle: {
    ...typography.h3,
    color: colors.text,
  },
  typeLabel: {
    ...typography.caption,
    color: colors.textSecondary,
    textTransform: "capitalize",
  },
  itemCard: {
    marginBottom: spacing.sm,
  },
  itemHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "flex-start",
    marginBottom: spacing.xs,
  },
  itemName: {
    ...typography.body,
    color: colors.text,
    fontWeight: "600",
    flex: 1,
    marginRight: spacing.sm,
  },
  itemDescription: {
    ...typography.small,
    color: colors.textSecondary,
    marginBottom: spacing.xs,
  },
  itemPrice: {
    ...typography.small,
    color: colors.primary,
    fontWeight: "600",
  },
  errorText: {
    ...typography.body,
    color: colors.error,
    marginBottom: spacing.md,
  },
  retryButton: {
    minWidth: 120,
  },
  emptyText: {
    ...typography.body,
    color: colors.textSecondary,
    textAlign: "center",
    paddingVertical: spacing.md,
  },
});
