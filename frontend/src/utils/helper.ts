import { ChainType } from "@/type";

import { always, cond, equals } from "ramda";


export type CoinInfo = {
  chainName: string;
  coinName: string;
  imgSrc: string;
  fees: number;
};

export const getChainNameChainType = cond([
  [
    equals<ChainType>("mgo"),
    always<CoinInfo>({
      chainName: "MGO",
      coinName: "MGO",
      imgSrc: "/images/mgo.png",
      fees: 0.3,
    }),
  ],
  [
    equals<ChainType>("sol"),
    always<CoinInfo>({
      chainName: "Solana",
      coinName: "Solana MGO",
      imgSrc: "/images/solana1.png",
      fees: 0.5,
    }),
  ],
]);
