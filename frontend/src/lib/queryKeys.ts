import { type TradeSymbol } from '@/types';

export const queryKeys = {
  user: {
    me: ['user'] as const,
    publicProfile: (username: string) => ['publicProfile', username] as const,
  },
  quotes: {
    detail: (symbol: TradeSymbol) => ['quote', symbol] as const,
    history: (symbol: TradeSymbol) => ['history', symbol] as const,
  },
  leaderboard: {
    all: ['leaderboard'] as const,
  },
  ladder: {
    active: ['ladder', 'active'] as const,
  },
} as const;
