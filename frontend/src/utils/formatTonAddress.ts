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

export function getAddressType(address: string) {
  const isMgo = /^0x[a-fA-F0-9]{64}$/.test(address);

  const isSolana = /^[1-9A-HJ-NP-Za-km-z]{32,44}$/.test(address);

  if (isMgo) return "mgo";
  if (isSolana) return "sol";

  return "unknown";
}
