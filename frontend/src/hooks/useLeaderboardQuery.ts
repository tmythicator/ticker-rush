import { useQuery } from '@tanstack/react-query';
import { getLeaderboard } from '@/lib/api';
import { queryKeys } from '@/lib/queryKeys';

export const useLeaderboardQuery = (limit: number = 10, offset: number = 0) => {
  return useQuery({
    queryKey: [...queryKeys.leaderboard.all, limit, offset],
    queryFn: () => getLeaderboard({ limit, offset }),
    refetchInterval: 60000,
  });
};
