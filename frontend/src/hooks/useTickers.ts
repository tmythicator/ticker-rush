import { useQuery } from '@tanstack/react-query';
import { getConfig } from '@/lib/api';

export const useTickers = () => {
  return useQuery({
    queryKey: ['tickers'],
    queryFn: getConfig,
    staleTime: Infinity, // Config rarely changes
    refetchOnWindowFocus: false,
  });
};
