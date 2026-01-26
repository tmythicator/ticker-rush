export const TradeAction = {
  BUY: 'BUY',
  SELL: 'SELL',
} as const;

export const TradeSymbol = {
  AAPL: 'AAPL',
  BTCUSDT: 'BINANCE:BTCUSDT',
} as const;

export const TradeSymbols = Object.values(TradeSymbol);
export type TradeSymbol = (typeof TradeSymbol)[keyof typeof TradeSymbol];
export type TradeAction = (typeof TradeAction)[keyof typeof TradeAction];

export const isTradeSymbol = (value: string): value is TradeSymbol =>
  (TradeSymbols as string[]).includes(value);
