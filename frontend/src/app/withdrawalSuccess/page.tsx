"use client";

import HeaderBar from "@/components/HeaderBar/HeaderBar";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { formatAddress } from "@/utils/formatTonAddress";
import { formatDateTime } from "@/utils/fromDate";
import { getChainNameChainType } from "@/utils/helper";
import Image from "next/image";
import { useRouter, useSearchParams } from "next/navigation";
import { path, pipe } from "ramda";
import { Suspense, useMemo } from "react";

function Withdrawal() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [props, chainBase] = useMemo(() => {
    const props = pipe(JSON.parse)(searchParams.get("props") ?? "");
    const chainBase = getChainNameChainType(path(["chainName"], props));
    return [props, chainBase];
  }, [searchParams]);

  return (
    <div className="h-full ">
      <HeaderBar title="Withdrawal successful" />
      <div className="h-[calc(100%-56px)] w-[93%] m-auto flex flex-col items-baseline">
        <div className="w-full flex flex-col gap-4  pt-28 pb-16 justify-center items-center">
          <Image src={"/images/success.png"} width={150} height={150} alt="" />
          <h5 className="text-white">Withdrawal successful</h5>
          <h6 className="text-white">
            - {path(["amount"], props)} {path(["coinName"], chainBase)}
          </h6>
        </div>
        <div className="w-full">
          <div
            className={cn(
              "w-full h-full gap-3 flex flex-col",
              "flex-1 text-[var(--text-s)] text-[12px] py-4 border-t border-[var(--t-border)]"
            )}
          >
            <div className="flex justify-between items-center">
              <p>TXID:</p>
              <p className="flex gap-2">
                {formatAddress(path(["data", "hash"], props))}
                <Image
                  src={path(["imgSrc"], chainBase)}
                  width={20}
                  height={20}
                  alt={""}
                />
              </p>
            </div>

            <div className="flex justify-between items-center">
              <p>To Addr:</p>
              <p>{formatAddress(path(["toAddress"], props))}</p>
            </div>

            <div className="flex justify-between items-center">
              <p>Date:</p>
              <p>{formatDateTime(path(["timestamp"], props))}</p>
            </div>
          </div>
        </div>

        <div className="w-full pt-14 flex justify-between items-end">
          <Button
            onClick={() => router.replace("/")}
            className="bg-[var(--t-main)] w-full"
          >
            Done
          </Button>
        </div>
      </div>
    </div>
  );
}

export default function WithdrawalSuccessPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Withdrawal />
    </Suspense>
  );
}
