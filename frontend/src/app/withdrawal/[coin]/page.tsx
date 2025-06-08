"use client";
import HeaderBar from "@/components/HeaderBar/HeaderBar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { balanceEffect, withdrawalEffect } from "@/effector/effector";
import { $balanceEffect } from "@/effector/store";
import { useWithdrawalEffect } from "@/hooks/useWithdrawalEffect";
import { cn } from "@/lib/utils";
import { ChainType } from "@/type";
import { WithdrawalReq } from "@/type/withdrawal";
import { getChainNameChainType } from "@/utils/helper";
import { useUnit } from "effector-react";
import { useParams } from "next/navigation";
import { mergeRight, path, pipe } from "ramda";
import { useEffect, useMemo } from "react";
import { useForm } from "react-hook-form";
import Image from "next/image";
import Smin from "@/components/Text/Smin";

export default function Withdrawal() {
  const { coin } = useParams<{
    coin: ChainType;
  }>();
  const chainBase = useMemo(() => getChainNameChainType(coin), [coin]);
  const { totalBalance } = useUnit($balanceEffect);

  const withdrawalPending = useUnit(withdrawalEffect.pending);

  const { handel } = useWithdrawalEffect();

  const { register, handleSubmit, watch } = useForm<WithdrawalReq>();

  const amountValue = watch("amount");

  useEffect(() => {
    balanceEffect({ chainName: coin });
  }, [coin]);

  return (
    <form
      onSubmit={handleSubmit(
        pipe((ev) => mergeRight(ev, { chainName: coin }), handel)
      )}
      className="h-full"
    >
      <div className="h-full ">
        <HeaderBar title={`${path(["coinName"], chainBase)} Withdrawal`} />
        <div className="h-[calc(100%-56px)] w-[93%] m-auto flex flex-col items-baseline justify-between">
          <div className="w-full flex flex-col gap-4 mt-4">
            <div>
              <Smin text="Receiving Address" />
              <div className="h-1"></div>
              <Textarea {...register("toAddress")} placeholder="" />
            </div>

            <div>
              <Smin text="Withdrawal Network" />
              <div className="h-1"></div>
              <div
                className={cn(
                  "text-white flex items-center py-3 px-2.5 rounded-[4px]",
                  "bg-[rgba(217,217,217,0.2)] gap-2 text-xs"
                )}
              >
                <div className="w-7 h-7 relative ">
                  <Image
                    src={path(["imgSrc"], chainBase)}
                    fill
                    className="object-contain"
                    alt={""}
                  />
                </div>
                <p>{path(["chainName"], chainBase)}</p>
              </div>
              {/* <Input
                {...register("chainName")}
                disabled
                className="uppercase"
                label=""
                value={coin}
              /> */}
            </div>

            <div>
              <Smin text="Withdrawal amount" />
              <div className="py-1">
                <Input
                  label=""
                  {...register("amount")}
                  placeholder=""
                />
              </div>

              <p className="text-xs text-[#999] flex justify-between text-right">
                <span>Available</span>
                <span>
                  {totalBalance} {path(["coinName"], chainBase)}
                </span>
              </p>
            </div>
          </div>
          <div className="flex-1 py-4 w-full">
            <div
              className={cn(
                "w-full h-full",
                "flex-1 text-[var(--text-s)] text-[12px] py-4 border-y border-[var(--t-border)]"
              )}
            >
              <p>1, Check the withdrawal address carefully to ensure it is correct.</p>
              <p>2, Confirm that the account balance is sufficient to meet the withdrawal requirements.</p>
            </div>
          </div>

          <div className="w-full py-4 flex justify-between items-end">
            <div>
              <p className="text-[10px] text-[var(--text-s)]">Amount of funds received</p>
              <p className="text-[12px] text-white">
                {amountValue ?? '--'} {path(["coinName"], chainBase)}
              </p>
              <p className="text-[10px] text-[var(--text-s)]">
                Network Fees{" "}
                <span className="text-white">
                  {path(["fees"], chainBase)} {path(["coinName"], chainBase)}
                </span>
              </p>
            </div>
            <div>
              <Button
                isLoading={withdrawalPending}
                className="bg-[var(--t-main)]"
              >
                Withdrawal
              </Button>
            </div>
          </div>
        </div>
      </div>
    </form>
  );
}
