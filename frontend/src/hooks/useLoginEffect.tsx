import { loginEffect, LoginProps } from "@/effector/effector";
import { useEffect } from "react";

export const useLoginEffect = () => {
  const handel = async (data: LoginProps) => {
    await loginEffect(data);
  };

  useEffect(() => {
    loginEffect.done.watch((payload) => {
      console.log("registerEffect doneData payload:", payload);
    });
  }, []);

  return { handel };
};
