import { useQuery } from '@tanstack/react-query';
import { getLeaderboard } from '@/lib/api';
import { QUERY_KEY_LEADERBOARD } from '@/lib/queryKeys';

export const useLeaderboardQuery = (limit: number = 10, offset: number = 0) => {
  return useQuery({
    queryKey: [...QUERY_KEY_LEADERBOARD, limit, offset],
    queryFn: () => getLeaderboard({ limit, offset }),
    refetchInterval: 60000,
  });
};
