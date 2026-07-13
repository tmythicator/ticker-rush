import { createTrade } from '@/lib/api';
import { queryKeys } from '@/lib/queryKeys';
import { TradeAction } from '@/types';
import type { User } from '@/types';
import { useMutation, useQueryClient } from '@tanstack/react-query';

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
      if (quantity <= 0) throw new Error('Quantity must be positive');

      return createTrade({ symbol: options.symbol, quantity, action });
    },
    onSuccess: (updatedProfile) => {
      queryClient.setQueryData<User>(queryKeys.user.me, (oldUser) => {
        if (!oldUser) return undefined;
        return {
          ...oldUser,
          balance: updatedProfile.balance,
          portfolio: updatedProfile.portfolio,
        };
      });
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
