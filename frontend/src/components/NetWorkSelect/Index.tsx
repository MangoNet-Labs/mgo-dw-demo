"use client";
import * as React from "react";
import { Check, ChevronDown } from "lucide-react";
import Image from "next/image";
import { cn } from "@/lib/utils";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { ChainType } from "@/type";
import { getChainNameChainType } from "@/utils/helper";
import { pathOr } from "ramda";
import { useState } from "react";
import { useRouter } from "next/navigation";

const frameworks: { value: ChainType; label: string }[] = [
  {
    value: "mgo",
    label: "Mango Network",
  },
  {
    value: "sol",
    label: "Solana",
  },
];

type ComboboxProps = {
  chain?: ChainType;
};
export function ComboboxDemo({ chain = "mgo" }: ComboboxProps) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState<ChainType>(chain);
  const chainBase = getChainNameChainType(chain);
  const { replace } = useRouter();

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <div className="bg-[rgba(217,217,217,.2)] w-32 justify-around flex px-3 py-2 rounded-3xl">
          <div className=" relative w-6 h-6">
            <Image
              src={pathOr("/images/solana1.png", ["imgSrc"], chainBase)}
              fill
              className="object-cover"
              alt={""}
            />
          </div>
          <p className="text-base text-white">
            {pathOr("mgo", ["chainName"], chainBase)}
          </p>
          <ChevronDown className="" color="white" />
        </div>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0">
        <Command>
          <CommandList>
            <CommandEmpty>No framework found.</CommandEmpty>
            <CommandGroup>
              {frameworks.map((framework) => (
                <CommandItem
                  key={framework.value}
                  value={framework.value}
                  onSelect={(currentValue) => {
                    setValue(currentValue as ChainType);
                    setOpen(false);
                    replace(`/${currentValue}`);
                  }}
                >
                  {framework.label}
                  <Check
                    className={cn(
                      "ml-auto",
                      value === framework.value ? "opacity-100" : "opacity-0"
                    )}
                  />
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
