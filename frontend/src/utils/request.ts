import { createEffect } from "effector";
import axios from "axios";
import { always, ifElse, isNil, path, pathEq, prop, uniq } from "ramda";
import { toast } from "react-toastify";
import { signoutEvent } from "@/services/events";

export interface RequestEffectPayload {
  url: string;
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
}

export const RequestBaseUrl = process.env.NEXT_PUBLIC_BASEURL;

export function createRequestEffect<
  R,
  P extends Record<string, unknown> | void
>({ url, method = "GET" }: RequestEffectPayload) {
  const keys = uniq(url.match(/:\w+/g) ?? []);
  let innerUrl = url;

  return createEffect<P extends void ? void : P & { _noAuth?: boolean }, R>(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    async (payload: any) => {
      const hasData = !["GET", "DELETE"].includes(method);
      const hasParams = !hasData && keys?.length === 0;

      if (keys?.length > 0 && payload != null) {
        innerUrl = keys.reduce(
          (acc, el) => acc.replaceAll(el, payload![el.slice(1)]),
          url
        );
      }

      const headers = { Token: "" };
      if (
        payload == null ||
        (typeof payload === "object" && !payload._noAuth)
      ) {
        const token = await getToken();
        if (token != null) {
          headers.Token = token;
        }
      }

      const res = await axios
        .request({
          url: RequestBaseUrl + innerUrl,
          method,
          headers,
          data: hasData ? payload : undefined,
          params: hasParams ? payload : undefined,
        })
        .then(prop("data"));

      if (pathEq(401, ["code"], res)) {
        signoutEvent();
      }
      if (res.code) {
        toast.error(res.message);
        throw res.message;
      }
      return res;
    }
  );
}

axios.interceptors.response.use(undefined, (error) => {
  const msg = path(["response", "data"], error) || error.message;
  toast.error(msg);

  if (pathEq(401, ["response", "status"], error)) {
    signoutEvent();
  }

  throw error;
});

export const getToken = () =>
  import("../effector/store").then(({ $user }) =>
    ifElse(
      isNil,
      always(undefined),
      (token) => `${token}`
    )($user.getState().token)
  );
