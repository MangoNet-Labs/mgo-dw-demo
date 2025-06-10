import * as React from "react";
import { cn } from "@/lib/utils";
import If from "../logic/If";
import { isNotEmpty } from "ramda";

type TextareaProps = React.ComponentProps<"textarea"> & {
  error?: string;
};

function Textarea({ className, error, ...props }: TextareaProps) {
  return (
    <div>
      <div className="bg-[var(--bg-color)] p-2 flex  rounded-[4px]">
        <textarea
          data-slot="textarea"
          className={cn(
            "text-white file:border-0 placeholder:text-muted-foreground  flex field-sizing-content min-h-16 w-full rounded-md px-3 py-2 text-base shadow-xs ",
            "appearance-none border-none outline-none bg-transparent resize-none p-0 m-0 shadow-none focus:ring-0",
            className
          )}
          {...props}
        />
      </div>
      <If prediction={isNotEmpty(error)}>
        <p className="text-red-600 text-xs mt-1">{error}</p>
      </If>
    </div>
  );
}

export { Textarea };
