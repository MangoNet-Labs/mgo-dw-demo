"use client";
import { ChevronLeft } from "lucide-react";
import { useRouter } from "next/navigation";
type Props = {
  title: string;
};
const HeaderBar = ({ title }: Props) => {
  const { back } = useRouter();
  return (
    <div className="relative flex items-center justify-center w-full h-14 px-4 bg-[rgba(41,49,54,0.60)]">
      <p onClick={back} className="absolute left-4">
        <ChevronLeft color="white" size={20} />
      </p>

      <p className="text-xl text-white">{title}</p>

      <div className="absolute right-4 w-5 h-5" />
    </div>
  );
};

export default HeaderBar;
