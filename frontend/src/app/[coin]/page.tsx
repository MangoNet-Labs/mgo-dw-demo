"use client";
import { HistoryCell } from "@/components/history/Cell";
import { ComboboxDemo } from "@/components/NetWorkSelect/Index";
import { Button } from "@/components/ui/button";
import { transactionEffect } from "@/effector/effector";
import { $transactionEffect, $user } from "@/effector/store";
import { ChainType } from "@/type";
import { copyText } from "@/utils/copy";
import { getChainNameChainType } from "@/utils/helper";
import { useUnit } from "effector-react";

import Link from "next/link";
import { useParams } from "next/navigation";
import { path } from "ramda";
import { useEffect, useMemo } from "react";

export default function CoinHome() {
  const params = useParams<{ coin?: ChainType }>();
  const coin = params.coin ?? "mgo";

  const { mgo_address, solana_address } = useUnit($user);

  const { list } = useUnit($transactionEffect);

  const address = useMemo(() => {
    if (coin == "mgo") return `${mgo_address}`;
    return `${solana_address}`;
  }, [mgo_address]);

  const ChainName = useMemo(() => getChainNameChainType(coin), [coin]);

  useEffect(() => {
    transactionEffect({
      page: 1,
      pageSize: 100,
      chainName: `${coin}`,
      type: 0,
    });
  }, [coin]);

  return (
    <div className="h-full w-[93%] pt-3.5 m-auto">
      <ComboboxDemo chain={coin} />
      <div className="w-full pt-6 text-center">
        <h4 className="text-base text-[var(--t-main)]">
          {path(["coinName"], ChainName)}
        </h4>
        <p className="text-[23px] text-center mt-2 text-white">
          Deposit and Withdrawal DEMO
        </p>
      </div>
      <div className="my-10 w-full">
        <Link href={`/withdrawal/${coin}`}>
          <Button className="bg-[var(--t-main)] w-full">
            Withdrawal {path(["coinName"], ChainName)}
          </Button>
        </Link>
      </div>
      <div className="flex text-white items-center justify-between">
        <p>{path(["coinName"], ChainName)} Address</p>
        {/* <Link
          href={`/qr-code?props=${JSON.stringify(
            mergeRight(ChainName, { address })
          )}`}
        >
          <ScanQrCode size={18} />
        </Link> */}
      </div>

      <div
        onClick={() => copyText(address)}
        className="mt-2 w-full p-1.5 rounded-xl bg-[var(--bg-color)] break-words text-white"
      >
        {address}
      </div>
      <div className="mt-3"></div>
      {list.map((ev) => (
        <HistoryCell {...ev} key={ev.hash} />
      ))}
    </div>
  );
}
