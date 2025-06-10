import BigNumber from "bignumber.js";
export function formatBalance(balance: string | null | undefined): string {
  if (!balance) return "--";

  return parseFloat(balance).toLocaleString("en-US", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 4,
  });
}

export function bgBalance(
  value: string | number | undefined,
  decimals: number = 9,
  fixed: number = 4
): string {
  if (!value) return "--";
  return new BigNumber(value)
    .dividedBy(new BigNumber(10).pow(decimals))
    .toFixed(fixed, BigNumber.ROUND_DOWN);
}
