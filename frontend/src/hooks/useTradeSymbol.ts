import { isTradeSymbol, type TradeSymbol } from '@/types';
import { useEffect, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useTickers } from './useTickers';

export const useTradeSymbol = () => {
  const [params, setParams] = useSearchParams();
  const { data, isLoading } = useTickers();
  const tickers = useMemo(() => data?.tickers ?? [], [data]);

  const rawSymbol = params.get('symbol');
  const isValid = rawSymbol && isTradeSymbol(rawSymbol, tickers);
  const symbol = (isValid ? rawSymbol : (tickers[0] ?? null)) as TradeSymbol | null;

  useEffect(() => {
    if (!isLoading && tickers.length > 0 && rawSymbol !== symbol) {
      setParams({ symbol: symbol! }, { replace: true });
    }
  }, [isLoading, tickers, rawSymbol, symbol, setParams]);

  return {
    symbol,
    setSymbol: (s: TradeSymbol) => setParams({ symbol: s }),
    tickers,
    isLoading,
  };
};
