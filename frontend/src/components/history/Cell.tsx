import { MgoTransaction } from "@/type/transaction";

type Props = {} & MgoTransaction;
export const HistoryCell = ({ hash, amount, chainName }: Props) => {
  return (
    <div className="w-full p-3 rounded-xl bg-[var(--bg-color)] mt-3">
      <h3 className="text-white text-base flex items-center justify-between mb-2">
        <span>Withdraw {chainName}</span>
        <span>
          {amount} {chainName}
        </span>
      </h3>
      <p className="text-[#999] text-xs">Trading Hours: 2024-10-17 14:32:18</p>
      <p className="text-[#999] text-xs">
        TXID: {hash}
      </p>

    </div>
  );
};
