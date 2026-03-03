/**
 * Formats a price (in VND minor units) as a human-readable Vietnamese Dong string.
 */
export function formatCurrency(amount: number): string {
  return `${amount.toLocaleString("vi-VN")}₫`;
}
