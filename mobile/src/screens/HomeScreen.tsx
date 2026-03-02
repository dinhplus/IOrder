import { View, Text, StyleSheet, ActivityIndicator } from "react-native";
import { StatusBar } from "expo-status-bar";
import { useHealth } from "@/hooks/useHealth";
import Button from "@/components/ui/Button";
import Card from "@/components/ui/Card";
import { colors, spacing, typography } from "@/styles/theme";

export default function HomeScreen(): React.JSX.Element {
  const { health, loading, error, refetch } = useHealth();

  return (
    <View style={styles.container}>
      <StatusBar style="auto" />

      <View style={styles.header}>
        <Text style={styles.title}>IOrder</Text>
        <Text style={styles.subtitle}>A modern restaurant ordering platform</Text>
      </View>

      <Card style={styles.statusCard}>
        <Text style={styles.cardTitle}>API Status</Text>

        {loading && (
          <View style={styles.centered}>
            <ActivityIndicator color={colors.primary} />
            <Text style={styles.statusText}>Checking connection…</Text>
          </View>
        )}

        {!loading && error && (
          <View style={styles.centered}>
            <Text style={[styles.statusIndicator, styles.statusError]}>●</Text>
            <Text style={styles.statusText}>Unable to reach the server</Text>
          </View>
        )}

        {!loading && !error && health && (
          <View>
            <View style={styles.row}>
              <Text style={styles.label}>Status</Text>
              <Text
                style={[
                  styles.value,
                  health.status === "ok" ? styles.statusOk : styles.statusError,
                ]}
              >
                {health.status}
              </Text>
            </View>
            <View style={styles.row}>
              <Text style={styles.label}>Database</Text>
              <Text
                style={[
                  styles.value,
                  health.db === "ok" ? styles.statusOk : styles.statusError,
                ]}
              >
                {health.db}
              </Text>
            </View>
            <View style={styles.row}>
              <Text style={styles.label}>Version</Text>
              <Text style={styles.value}>{health.version}</Text>
            </View>
          </View>
        )}
      </Card>

      <Button
        title="Refresh"
        onPress={refetch}
        loading={loading}
        variant="outline"
        style={styles.refreshButton}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: spacing.lg,
    backgroundColor: colors.surface,
  },
  header: {
    alignItems: "center",
    marginTop: spacing.xxl,
    marginBottom: spacing.xl,
  },
  title: {
    ...typography.h1,
    color: colors.primary,
    marginBottom: spacing.sm,
  },
  subtitle: {
    ...typography.body,
    color: colors.textSecondary,
    textAlign: "center",
  },
  statusCard: {
    marginBottom: spacing.lg,
  },
  cardTitle: {
    ...typography.h3,
    color: colors.text,
    marginBottom: spacing.md,
  },
  centered: {
    alignItems: "center",
    paddingVertical: spacing.md,
  },
  statusIndicator: {
    fontSize: 24,
    marginBottom: spacing.xs,
  },
  statusText: {
    ...typography.body,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  row: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingVertical: spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  label: {
    ...typography.body,
    color: colors.textSecondary,
  },
  value: {
    ...typography.body,
    color: colors.text,
    fontWeight: "600",
  },
  statusOk: {
    color: colors.success,
  },
  statusError: {
    color: colors.error,
  },
  refreshButton: {
    alignSelf: "center",
    minWidth: 140,
  },
});
