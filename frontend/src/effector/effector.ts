import { GetBalanceByCoinReq, GetBalanceByCoinResp } from "@/type/balance";
import { TransactionListReq, TransactionListResp } from "@/type/transaction";
import { SigninData } from "@/type/user";
import { WithdrawalReq, WithdrawalResp } from "@/type/withdrawal";
import { createRequestEffect } from "@/utils/request";

export type LoginProps = {
  username: string;
  password: string;
};

export type RegisterProps = {
  rePassword: string;
} & LoginProps;

export const registerEffect = createRequestEffect<
  Partial<SigninData>,
  LoginProps
>({
  url: "register",
  method: "POST",
});

export const loginEffect = createRequestEffect<Partial<SigninData>, LoginProps>(
  {
    url: "login",
    method: "POST",
  }
);

export const transactionEffect = createRequestEffect<
  TransactionListResp,
  TransactionListReq
>({
  url: "transaction",
  method: "POST",
});


export const balanceEffect  = createRequestEffect<
  GetBalanceByCoinResp,
  GetBalanceByCoinReq
>({
  url: "balance",
  method: "POST",
});




export const withdrawalEffect  = createRequestEffect<
  WithdrawalResp,
  WithdrawalReq
>({
  url: "withdrawal",
  method: "POST",
});
