import { getHistory } from '@/lib/api';
import { QUERY_KEY_HISTORY } from '@/lib/queryKeys';
import { type TradeSymbol } from '@/types';
import { useQuery } from '@tanstack/react-query';
import { type ISeriesApi, type Time } from 'lightweight-charts';
import { useEffect } from 'react';

export const useChartHistory = (
  symbol: TradeSymbol | null,
  seriesRef: React.RefObject<ISeriesApi<'Area'> | null>,
) => {
  const { data: history } = useQuery({
    queryKey: QUERY_KEY_HISTORY(symbol || ''),
    queryFn: () => getHistory(symbol!, 100),
    enabled: !!symbol,
    staleTime: 1000 * 60 * 5,
    select: (data) => {
      return data.map((q) => ({
        time: q.timestamp as Time,
        value: q.price,
      }));
    },
  });

  useEffect(() => {
    if (!symbol || !seriesRef.current) return;

    if (!history) {
      seriesRef.current.setData([]);
      return;
    }

    try {
      seriesRef.current.setData(history);
    } catch (err) {
      console.error('Failed to set chart data:', err);
    }
  }, [symbol, history, seriesRef]);
};
