import { getHistory } from '@/lib/api';
import { queryKeys } from '@/lib/queryKeys';
import { queryConfig } from '@/lib/queryConfig';
import { type TradeSymbol } from '@/types';
import { useQuery } from '@tanstack/react-query';
import { type ISeriesApi, type Time } from 'lightweight-charts';
import { useEffect } from 'react';

export const useChartHistory = (
  symbol: TradeSymbol | null,
  seriesRef: React.RefObject<ISeriesApi<'Area'> | null>,
) => {
  const { data: history } = useQuery({
    queryKey: queryKeys.quotes.history(symbol || ''),
    queryFn: () => getHistory({ symbol: symbol!, limit: 1000 }),
    enabled: !!symbol,
    ...queryConfig.history,
    select: (data) => {
      return data.map((q) => ({
        time: (q.timestamp ? Math.floor(q.timestamp.getTime() / 1000) : 0) as Time,
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
