
export type WithdrawalReq = {
  toAddress: string;
  chainName: string;
  amount: number;
};

export interface WithdrawalResp {
  toAddress: string;
  chainName: string;
  amount: number;
}
