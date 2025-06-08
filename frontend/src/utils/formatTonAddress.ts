export function formatAddress(address: string | null | undefined): string {
  if (!address) return "--";
  if (!address || address.length < 8) {
    return "Invalid address";
  }
  try {
    const normalizedAddress = address;
    const firstFour = normalizedAddress.slice(0, 4);
    const lastFour = normalizedAddress.slice(-4);
    return `${firstFour}...${lastFour}`;
  } catch (err) {
    console.error(err);
    return "err";
  }
}
