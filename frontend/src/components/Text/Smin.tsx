import { cn } from "@/lib/utils";

type Props = {
  text: string;
  className?: string;
};

const SminText = ({ text, className }: Props) => {
  return <p className={cn("text-xs text-white", className)}>{text}</p>;
};

export default SminText;
