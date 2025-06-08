import { cn } from "@/lib/utils";
import Link from "next/link";

type Props = {
  type: string;
  acType: string;
  coin: string;
};

export const TrType = ({ coin, type, acType }: Props) => {
  return (
    <Link
      href={`/${coin}?type=${type}`}
      onClick={() => {}}
      className={cn(
        "rounded-4xl border-[1px] border-[var(--text-s)] capitalize px-3 py-1 text-xs text-[var(--text-s)]",
        {
          "bg-[var(--text-s)] text-white": type == acType,
        }
      )}
    >
      {type}
    </Link>
  );
};
