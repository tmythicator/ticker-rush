import type { TickerInfo } from './lib/proto/ladder/v1/ladder';

export * from './lib/api';

export type { Quote } from './lib/proto/exchange/v1/exchange';
export { GetLeaderboardResponse } from './lib/proto/leaderboard/v1/leaderboard';
export type { LeaderboardEntry } from './lib/proto/leaderboard/v1/leaderboard';
export type { PortfolioItem } from './lib/proto/portfolio/v1/portfolio';
export type { UpdateUserRequest, User } from './lib/proto/user/v1/user';
export type { TickerInfo };

export const TradeAction = {
  BUY: 'BUY',
  SELL: 'SELL',
} as const;

export type TradeSymbol = string;
export type TradeAction = (typeof TradeAction)[keyof typeof TradeAction];

export const isTradeSymbol = (value: string, validTickers: TickerInfo[]): value is TradeSymbol =>
  validTickers.some((t) => t.symbol === value);

export type TickerSource = 'Finnhub' | 'CoinGecko';
