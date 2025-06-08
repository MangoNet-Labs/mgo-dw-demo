import {
  clientSideRenderredEvent,
  clientSideRenderringEvent,
} from "@/services/events";
import { createStore } from "effector";

export function getSavedState(key: string) {
  const saved = localStorage.getItem(key);

  return saved == null ? {} : JSON.parse(saved);
}

export default function createStorage<T extends object>(
  key: string,
  initial: T
) {
  key = `dw-demo-${key}`;

  const store = createStore(initial).on(clientSideRenderringEvent, () =>
    getSavedState(key)
  );

  const unwatch = clientSideRenderringEvent.watch(() => {
    unwatch();
    store.updates.watch((state) => {
      globalThis.localStorage?.setItem(key, JSON.stringify(state));
    });

    setTimeout(clientSideRenderredEvent);
  });

  return store;
}
