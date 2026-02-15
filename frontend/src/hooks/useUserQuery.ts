import { useQuery } from '@tanstack/react-query';
import { getUser } from '@/lib/api';
import { QUERY_KEY_USER } from '@/lib/queryKeys';

export const useUserQuery = () => {
  return useQuery({
    queryKey: QUERY_KEY_USER,
    queryFn: getUser,
    retry: false,
    staleTime: 5 * 60 * 1000, // Keep user data fresh for 5 mins
  });
};
