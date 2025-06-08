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

  id: string;
  hash: string;
  chainName: string;
  type: number;
  timestamp: string;
  amount: string;
  status: string;
};
