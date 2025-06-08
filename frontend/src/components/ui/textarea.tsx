import * as React from "react";
import { cn } from "@/lib/utils";

function Textarea({ className, ...props }: React.ComponentProps<"textarea">) {
  return (
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
  );
}

export { Textarea };
