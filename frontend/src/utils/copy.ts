import { toast } from "react-toastify";

export const copyText = async (text: string, tis?: string) => {
  await navigator.clipboard.writeText(text);
  if (tis) {
    return toast.success(`${tis}!`);
  }
  toast.success("Success!");
};
