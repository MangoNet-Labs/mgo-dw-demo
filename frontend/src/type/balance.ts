
export type GetBalanceByCoinReq = {
  chainName: string;
};
export interface GetBalanceByCoinResp {
  coinType: string;
  coinObjectCount: number;
  totalBalance: string;
}
