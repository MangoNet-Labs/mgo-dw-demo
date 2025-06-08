import { withdrawalEffect } from "@/effector/effector";
import { WithdrawalReq } from "@/type/withdrawal";
import { useRouter } from "next/navigation";

import { evolve, mergeRight } from "ramda";
import { useEffect } from "react";

export const useWithdrawalEffect = () => {
  const { push } = useRouter();
  const handel = async (data: WithdrawalReq) => {
    await withdrawalEffect(
      evolve({
        amount: Number,
      })(data)
    );
  };

  useEffect(() => {
    withdrawalEffect.done.watch(({ params, result }) => {
      push(
        `/withdrawalSuccess?props=${JSON.stringify(mergeRight(params, result))}`
      );
      return;
    });
  }, []);

  return { handel };
};
