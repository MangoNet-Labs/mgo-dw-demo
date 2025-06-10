import { registerEffect, RegisterProps } from "@/effector/effector";
import { toast } from "react-toastify";

export const useRegisterEffect = () => {
  const handel = async (data: RegisterProps) => {
    if (data.password != data.rePassword) {
      toast.error("Passwords do not match");
      return;
    }
    await registerEffect(data);
  };
  return { handel };
};
