import { useState } from 'react';
import { buyStock, sellStock } from '../lib/api';
import { useAuth } from './useAuth';

import { TradeAction } from '../types';

interface UseTradeOptions {
  symbol: string;
  onSuccess?: () => void;
}

export const useTrade = (options: UseTradeOptions) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { refreshUser } = useAuth();

  const executeTrade = async (action: TradeAction, quantity: number) => {
    if (!quantity || quantity <= 0) return;

    setIsLoading(true);
    setError(null);

    try {
      let updatedUser;
      if (action === TradeAction.BUY) {
        updatedUser = await buyStock(options.symbol, quantity);
      } else {
        updatedUser = await sellStock(options.symbol, quantity);
      }
      refreshUser(updatedUser);
      options.onSuccess?.();
    } catch (e) {
      if (e instanceof Error) {
        setError(e.message);
      } else {
        setError('Trade failed');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return {
    executeTrade,
    isLoading,
    error,
    setError,
  };
};
