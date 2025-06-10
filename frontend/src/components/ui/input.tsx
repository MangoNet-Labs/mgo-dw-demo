"use client";
import * as React from "react";
import { cn } from "@/lib/utils";
import If from "../logic/If";
import { EyeClosed, Eye } from "lucide-react";
import { useState } from "react";
import IfElse from "../logic/if-else";
import { isNotEmpty } from "ramda";

type InputProps = React.ComponentProps<"input"> & {
  label: string;
  error?: string;
};

function Input({ className, type, label, error, ...props }: InputProps) {
  const [isShow, setShow] = useState(false);
  return (
    <div>
      <div
        className={cn(
          "text-white flex items-center py-3 px-2.5 rounded-[4px]",
          "bg-[rgba(217,217,217,0.2)] gap-2 text-xs",
          className
        )}
      >
        <If prediction={isNotEmpty(label)}>
          <p>{label}</p>
        </If>

        <input
          type={isShow ? "text" : type}
          data-slot="input"
          className={cn(
            "selection:text-primary-foreground  border-input  shadow-xs transition-[color,box-shadow] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none",
            "flex-1",
            "aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"
          )}
          {...props}
        />
        <If prediction={type == "password" || isShow == true}>
          <div className="font-[14px]" onClick={() => setShow(!isShow)}>
            <IfElse prediction={isShow}>
              <Eye size={14} />
              <EyeClosed size={14} />
            </IfElse>
          </div>
        </If>
      </div>
      <If prediction={isNotEmpty(error)}>
        <p className="text-red-600 text-xs mt-1">{error}</p>
      </If>
    </div>
  );
}

export { Input };
