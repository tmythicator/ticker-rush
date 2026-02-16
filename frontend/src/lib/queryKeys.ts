import { type TradeSymbol } from '@/types';

export const QUERY_KEY_USER = ['user'];
export const QUERY_KEY_QUOTE = (symbol: TradeSymbol) => ['quote', symbol];
export const QUERY_KEY_HISTORY = (symbol: TradeSymbol) => ['history', symbol];
