"use client";
import { HistoryCell } from "@/components/history/Cell";
import IfElse from "@/components/logic/if-else";
import { ComboboxDemo } from "@/components/NetWorkSelect/Index";
import { Button } from "@/components/ui/button";
import { transactionEffect } from "@/effector/effector";
import { $transactionEffect, $user } from "@/effector/store";
import { ChainType } from "@/type";
import { copyText } from "@/utils/copy";
import { getChainNameChainType } from "@/utils/helper";
import { useUnit } from "effector-react";
import Link from "next/link";
import { useParams, useSearchParams } from "next/navigation";
import { findIndex, mergeRight, path } from "ramda";
import { Suspense, useEffect, useMemo } from "react";
import Image from "next/image";
import { TrType } from "@/components/Tags/TrType";
import dynamic from "next/dynamic";
const btns = ["all", "deposit", "withdraw"];

function CoinHome() {
  const searchParams = useSearchParams();

  const type = searchParams.get("type") ?? "all";
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
      type: findIndex((e) => e == type, btns),
    });
  }, [coin, type]);

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
      <div className="my-5 w-full">
        <Link href={`/withdrawal/${coin}`}>
          <Button className="bg-[var(--t-main)] w-full">
            Withdrawal {path(["coinName"], ChainName)}
          </Button>
        </Link>
      </div>
      <div className="flex text-white items-center justify-between">
        <p>{path(["coinName"], ChainName)} Address</p>
      </div>

      <div
        onClick={() => copyText(address)}
        className="mt-2 w-full p-1.5 rounded-[4px] bg-[var(--bg-color)] break-words text-white"
      >
        {address}
      </div>
      <div className="mt-6 flex gap-3">
        {btns.map((e) => (
          <TrType coin={coin} type={e} acType={type} key={e} />
        ))}
      </div>

      <IfElse prediction={list.length > 0}>
        <>
          {list.map((ev) => (
            <HistoryCell
              {...mergeRight(ev, ChainName)}
              address={address}
              key={ev.digest}
            />
          ))}
        </>
        <div className="text-center pt-10">
          <Image
            src={"/images/niu_data.png"}
            width={150}
            height={50}
            alt={""}
            className=" m-auto"
          />
          <p className="text-[#999] text-base mt-5">No more data yet~</p>
        </div>
      </IfElse>
    </div>
  );
}

const DynamicCoinHome = dynamic(() => Promise.resolve(CoinHome), {
  ssr: false,
});

export default function CoinHomePage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <DynamicCoinHome />
    </Suspense>
  );
}
