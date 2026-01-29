import { useMutation, useQueryClient } from '@tanstack/react-query';
import { buyStock, sellStock } from '../lib/api';
import { QUERY_KEY_USER } from '../lib/queryKeys';
import { TradeAction } from '../types';

interface UseTradeOptions {
  symbol: string;
  onSuccess?: () => void;
}

interface TradeMutation {
  action: TradeAction;
  quantity: number;
}

export const useTrade = (options: UseTradeOptions) => {
  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn: ({ action, quantity }: TradeMutation) => {
      if (quantity <= 0) throw new Error('Qunatity must be positive');

      return action === TradeAction.BUY
        ? buyStock(options.symbol, quantity)
        : sellStock(options.symbol, quantity);
    },
    onSuccess: (updatedUser) => {
      queryClient.setQueryData(QUERY_KEY_USER, updatedUser);
      options.onSuccess?.();
    },
  });

  return {
    executeTrade: mutation.mutate,
    isLoading: mutation.isPending,
    error: mutation.error?.message ?? null,
    setError: () => mutation.reset(),
  };
};
