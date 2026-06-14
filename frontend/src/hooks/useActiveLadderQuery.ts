import { useQuery } from '@tanstack/react-query';
import { getActiveLadder } from '@/lib/api';
import { queryKeys } from '@/lib/queryKeys';

export const useActiveLadderQuery = () => {
  return useQuery({
    queryKey: queryKeys.ladder.active,
    queryFn: getActiveLadder,
  });
};
