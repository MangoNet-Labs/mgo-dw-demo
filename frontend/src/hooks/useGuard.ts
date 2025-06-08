import { useEffect } from "react";
import { clientSideRenderredEvent } from "../services/events";
import { $user } from "@/effector/store";
import { usePathname, useRouter } from "next/navigation";

const unsecuredPaths = ["/login", "/register"];

export default function useGuard() {
  const pathname = usePathname();
  const { push } = useRouter();

  useEffect(
    () =>
      clientSideRenderredEvent.watch(() => {
        if (
          $user.getState().token == null &&
          !unsecuredPaths.includes(pathname)
        ) {
          push("/login");
        }
      }),
    [pathname]
  );
}
