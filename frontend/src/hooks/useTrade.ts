import { useState } from 'react';
import { buyStock, sellStock } from '../lib/api';

import { TradeAction } from '../types';

interface UseTradeOptions {
    userId: number;
    symbol: string;
    onSuccess?: () => void;
}

export const useTrade = (options: UseTradeOptions) => {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const executeTrade = async (action: TradeAction, quantity: number) => {
        if (!quantity || quantity <= 0) return;

        setIsLoading(true);
        setError(null);

        try {
            if (action === TradeAction.BUY) {
                await buyStock(options.userId, options.symbol, quantity);
            } else {
                await sellStock(options.userId, options.symbol, quantity);
            }
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
        setError
    };
};
