import type { TickerInfo } from './lib/proto/ladder/v1/ladder';
export type { Quote } from './lib/proto/exchange/v1/exchange';
export { GetLeaderboardResponse } from './lib/proto/leaderboard/v1/leaderboard';
export type { LeaderboardEntry } from './lib/proto/leaderboard/v1/leaderboard';
export type { PortfolioItem, UpdateUserRequest, User } from './lib/proto/user/v1/user';
export type { Ladder } from './lib/proto/ladder/v1/ladder';
export type { TickerInfo };

import { TradeAction as ApiTradeAction } from './lib/proto/exchange/v1/exchange';

export const TradeAction = {
  BUY: ApiTradeAction.TRADE_ACTION_BUY,
  SELL: ApiTradeAction.TRADE_ACTION_SELL,
} as const;

export type TradeSymbol = string;
export type TradeAction = (typeof TradeAction)[keyof typeof TradeAction];

export const isTradeSymbol = (value: string, validTickers: TickerInfo[]): value is TradeSymbol =>
  validTickers.some((t) => t.symbol === value);

export type TickerSource = 'Finnhub' | 'CoinGecko' | 'CG' | 'FH';
