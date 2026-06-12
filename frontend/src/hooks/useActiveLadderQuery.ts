import { useQuery } from '@tanstack/react-query';
import { getActiveLadder } from '@/lib/api';
import { QUERY_KEY_ACTIVE_LADDER } from '@/lib/queryKeys';

export const useActiveLadderQuery = () => {
  return useQuery({
    queryKey: QUERY_KEY_ACTIVE_LADDER,
    queryFn: getActiveLadder,
  });
};
