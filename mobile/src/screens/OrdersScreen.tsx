import { View, Text, StyleSheet, ActivityIndicator, ScrollView } from "react-native";
import { StatusBar } from "expo-status-bar";
import { useOrders } from "@/hooks/useOrders";
import Card from "@/components/ui/Card";
import Badge, { type BadgeVariant } from "@/components/ui/Badge";
import Button from "@/components/ui/Button";
import { formatCurrency } from "@/lib/format";
import { colors, spacing, typography } from "@/styles/theme";
import { ORDER_STATUS_DISPLAY, type OrderStatus } from "@/types/order";

function statusVariant(status: OrderStatus): BadgeVariant {
  switch (status) {
    case "READY":
    case "SERVED":
    case "PAID":
    case "CLOSED":
      return "success";
    case "CANCELLED":
    case "REJECTED":
      return "error";
    case "IN_PREPARATION":
      return "info";
    case "PAYMENT_REQUESTED":
      return "warning";
    default:
      return "default";
  }
}

export default function OrdersScreen(): React.JSX.Element {
  const { orders, loading, error, refetch } = useOrders();

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
        <Text style={styles.errorText}>Unable to load orders</Text>
        <Button title="Retry" onPress={refetch} style={styles.retryButton} />
      </View>
    );
  }

  type OrdersMap = Record<OrderStatus, typeof orders>;
  const grouped = ORDER_STATUS_DISPLAY.reduce<OrdersMap>(
    (acc, s) => {
      acc[s] = orders.filter((o) => o.status === s);
      return acc;
    },
    {} as OrdersMap,
  );

  const activeStatuses = ORDER_STATUS_DISPLAY.filter((s) => grouped[s].length > 0);

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <StatusBar style="auto" />

      <View style={styles.header}>
        <Text style={styles.title}>Orders</Text>
        <Text style={styles.subtitle}>{orders.length} total</Text>
      </View>

      {orders.length === 0 ? (
        <Text style={styles.emptyText}>No orders found.</Text>
      ) : (
        activeStatuses.map((status) => (
          <View key={status} style={styles.section}>
            <View style={styles.sectionHeader}>
              <Badge
                label={status.replace(/_/g, " ")}
                variant={statusVariant(status)}
              />
              <Text style={styles.countLabel}>
                {grouped[status].length} order{grouped[status].length !== 1 ? "s" : ""}
              </Text>
            </View>

            {grouped[status].map((order) => (
              <Card key={order.id} style={styles.orderCard}>
                <View style={styles.orderHeader}>
                  <Text style={styles.orderId}>#{order.id.slice(0, 8)}</Text>
                  <Text style={styles.orderTable}>Table: {order.table_id.slice(0, 8)}</Text>
                </View>
                <View style={styles.row}>
                  <Text style={styles.label}>Total</Text>
                  <Text style={styles.value}>{formatCurrency(order.total)}</Text>
                </View>
                <View style={styles.row}>
                  <Text style={styles.label}>Created</Text>
                  <Text style={styles.value}>
                    {new Date(order.created_at).toLocaleString("vi-VN")}
                  </Text>
                </View>
                {order.notes ? (
                  <View style={styles.row}>
                    <Text style={styles.label}>Notes</Text>
                    <Text style={[styles.value, styles.noteValue]} numberOfLines={2}>
                      {order.notes}
                    </Text>
                  </View>
                ) : null}
              </Card>
            ))}
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
    gap: spacing.sm,
    marginBottom: spacing.md,
  },
  countLabel: {
    ...typography.small,
    color: colors.textSecondary,
  },
  orderCard: {
    marginBottom: spacing.sm,
  },
  orderHeader: {
    marginBottom: spacing.sm,
  },
  orderId: {
    ...typography.body,
    color: colors.text,
    fontWeight: "600",
  },
  orderTable: {
    ...typography.small,
    color: colors.textSecondary,
  },
  row: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingVertical: spacing.xs,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  label: {
    ...typography.small,
    color: colors.textSecondary,
  },
  value: {
    ...typography.small,
    color: colors.text,
    fontWeight: "500",
  },
  noteValue: {
    maxWidth: 160,
    textAlign: "right",
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
