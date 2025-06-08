export function formatBalance(balance: string | null | undefined): string {
  if (!balance) return "--";

  return parseFloat(balance).toLocaleString("en-US", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 3,
  });
}
