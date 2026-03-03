import { View, Text, StyleSheet, ActivityIndicator, ScrollView } from "react-native";
import { StatusBar } from "expo-status-bar";
import { useTables } from "@/hooks/useTables";
import Card from "@/components/ui/Card";
import Badge, { type BadgeVariant } from "@/components/ui/Badge";
import Button from "@/components/ui/Button";
import { colors, spacing, typography } from "@/styles/theme";

function tableStatusVariant(status: string): BadgeVariant {
  if (status === "available") return "success";
  if (status === "occupied") return "error";
  if (status === "reserved") return "warning";
  return "default";
}

export default function TablesScreen(): React.JSX.Element {
  const { data, loading, error, refetch } = useTables();

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
        <Text style={styles.errorText}>Unable to load floor plans</Text>
        <Button title="Retry" onPress={refetch} style={styles.retryButton} />
      </View>
    );
  }

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <StatusBar style="auto" />

      <View style={styles.header}>
        <Text style={styles.title}>Floor Plans</Text>
        <Text style={styles.subtitle}>
          {data.length} floor plan{data.length !== 1 ? "s" : ""}
        </Text>
      </View>

      {data.length === 0 ? (
        <Text style={styles.emptyText}>No floor plans found.</Text>
      ) : (
        data.map(({ plan, tables }) => (
          <View key={plan.id} style={styles.section}>
            <View style={styles.sectionHeader}>
              <Text style={styles.sectionTitle}>{plan.name}</Text>
              <Text style={styles.floorLabel}>Floor {plan.floor_level}</Text>
              <Badge
                label={plan.is_active ? "Active" : "Inactive"}
                variant={plan.is_active ? "success" : "default"}
              />
              <Text style={styles.tableCount}>{tables.length} tables</Text>
            </View>

            {tables.length === 0 ? (
              <Text style={styles.emptyText}>No tables on this floor plan.</Text>
            ) : (
              <View style={styles.tableGrid}>
                {tables.map((table) => (
                  <Card key={table.id} style={styles.tableCard}>
                    <Text style={styles.tableName}>{table.name}</Text>
                    <Text style={styles.tableCapacity}>{table.capacity} seats</Text>
                    <Badge
                      label={table.status}
                      variant={tableStatusVariant(table.status)}
                    />
                  </Card>
                ))}
              </View>
            )}
          </View>
        ))
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
  floorLabel: {
    ...typography.small,
    color: colors.textSecondary,
  },
  tableCount: {
    ...typography.small,
    color: colors.textSecondary,
  },
  tableGrid: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
  },
  tableCard: {
    width: "47%",
    alignItems: "center",
    paddingVertical: spacing.md,
  },
  tableName: {
    ...typography.body,
    color: colors.text,
    fontWeight: "600",
    marginBottom: spacing.xs,
  },
  tableCapacity: {
    ...typography.caption,
    color: colors.textSecondary,
    marginBottom: spacing.xs,
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
