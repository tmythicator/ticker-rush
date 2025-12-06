export const TradeAction = {
    BUY: "BUY",
    SELL: "SELL"
} as const;

export type TradeAction = typeof TradeAction[keyof typeof TradeAction];
