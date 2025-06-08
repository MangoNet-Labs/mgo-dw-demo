import createStorage from "@/utils/storage";
import {
  balanceEffect,
  loginEffect,
  registerEffect,
  transactionEffect,
} from "./effector";
import { signoutEvent } from "@/services/events";
import { SigninData } from "@/type/user";
import { asPayload } from "@/utils/reducer";
import { createStore } from "effector";
import { TransactionListResp } from "@/type/transaction";
import { GetBalanceByCoinResp } from "@/type/balance";

export const $user = createStorage<Partial<SigninData>>("user", {})
  .on(registerEffect.doneData, asPayload)
  .on(loginEffect.doneData, asPayload)
  .on(signoutEvent, () => ({}));

export const $transactionEffect = createStore<TransactionListResp>({
  page: 0,
  pageSize: 0,
  total: 0,
  list: [],
}).on(transactionEffect.doneData, asPayload);

export const $balanceEffect = createStore<GetBalanceByCoinResp>({
  coinType: "",
  coinObjectCount: 0,
  totalBalance: "",
}).on(balanceEffect.doneData, asPayload);
