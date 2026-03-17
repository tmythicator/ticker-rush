import { getActiveLadder } from '@/lib/api';
import { useQuery } from '@tanstack/react-query';

export const useTickers = () => {
  const {
    data: ladder,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['ladder', 'active'],
    queryFn: getActiveLadder,
    staleTime: Infinity,
    refetchOnWindowFocus: false,
  });

  return {
    data: ladder?.allowed_tickers,
    isLoading,
    error,
  };
};
