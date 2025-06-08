import { withdrawalEffect } from "@/effector/effector";
import { WithdrawalReq } from "@/type/withdrawal";
import { useRouter } from "next/navigation";

import { evolve } from "ramda";
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
    withdrawalEffect.done.watch((payload) => {
      push(`/withdrawalSuccess?props=${JSON.stringify(payload)}`);
      return;
    });
  }, []);

  return { handel };
};
