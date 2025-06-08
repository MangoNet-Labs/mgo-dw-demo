import { path } from "ramda";

export function asPayload<T>(_: T, payload: T) {

  return path<T>(["data"], payload);
}

export function mergePayload<T>(state: T, payload: Partial<T>) {
  return { ...state, ...payload };
}
