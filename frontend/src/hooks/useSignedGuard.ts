import { $user } from "@/effector/store";
import { useUnit } from "effector-react";
import { useRouter } from "next/navigation";

import { useEffect } from "react";

export default function useSignedGuard() {
  const userState = useUnit($user);
  const router = useRouter();

  useEffect(() => {
    if (userState.token != null) {
      router.replace("/mgo");
    }
  }, [router, userState.token]);
}
