export * from './lib/api';

export type { Quote } from './lib/proto/exchange/v1/exchange';
export { GetLeaderboardResponse } from './lib/proto/leaderboard/v1/leaderboard';
export type { LeaderboardEntry } from './lib/proto/leaderboard/v1/leaderboard';
export type { PortfolioItem, User } from './lib/proto/user/v1/user';

export const TradeAction = {
  BUY: 'BUY',
  SELL: 'SELL',
} as const;

export type TradeSymbol = string;
export type TradeAction = (typeof TradeAction)[keyof typeof TradeAction];

export const isTradeSymbol = (value: string, validTickers: string[]): value is TradeSymbol =>
  validTickers.includes(value);

export type TickerSource = 'FH' | 'CG';
