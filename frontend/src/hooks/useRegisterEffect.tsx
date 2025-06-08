import { LoginProps, registerEffect } from "@/effector/effector";

export const useRegisterEffect = () => {
  const handel = async (data: LoginProps) => {
    await registerEffect(data);
  };
  return { handel };
};
