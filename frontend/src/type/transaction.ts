export type TransactionListReq = {
  page: number;
  pageSize: number;
  chainName: string;
  type: number;
};

export type TransactionListResp = {
  page: number;
  pageSize: number;
  total: number;
  list: MgoTransaction[];
};

export type MgoTransaction = {
  id: number;
  digest: string;
  from: string;
  to: string;
  amount?: string;
  from_amount?: string;
  checkpoint?: string;
  coin_type?: string;
  gas_owner?: string;
  gas_price?: string;
  gas_budget?: string;
  timestamp_ms: string;
};
