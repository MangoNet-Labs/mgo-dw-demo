import { MgoTransaction } from "@/type/transaction";
import { formatAddress } from "@/utils/formatTonAddress";
import { formatDateTime } from "@/utils/fromDate";
import { CoinInfo } from "@/utils/helper";
import { always, ifElse } from "ramda";
import { useMemo } from "react";
import Image from "next/image";

type Props = {
  address: string;
} & MgoTransaction &
  CoinInfo;
export const HistoryCell = ({
  amount,
  coinName,
  from,
  chainName,
  fees,
  timestamp_ms,
  address,
  imgSrc,
}: Props) => {
  const transferType = useMemo(
    ifElse(() => address == from, always("Withdraw"), always("Recharge")),
    []
  );
  return (
    <div className="w-full p-3 rounded-[4px] bg-[var(--bg-color)] mt-3">
      <h3 className="text-white text-base flex items-center justify-between mb-2">
        <span>
          {transferType} {coinName}
        </span>
        <span>
          {amount} {coinName}
        </span>
      </h3>
      <div className="flex flex-col gap-1">
        <p className="text-[#999] text-xs">
          Date: {formatDateTime(timestamp_ms)}
        </p>
        <div className="flex justify-between items-center">
          <p className="text-[#999] text-xs">
            From Addr: {formatAddress(from)}
          </p>
          <Image src={imgSrc} width={20} height={20} alt={""} />
        </div>
        <p className="text-[#999] text-xs">
          Fees: {fees}
          {chainName}
        </p>
      </div>
    </div>
  );
};
