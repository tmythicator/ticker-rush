import { useQuery } from '@tanstack/react-query';
import { getUser } from '@/lib/api';
import { queryKeys } from '@/lib/queryKeys';

export const useUserQuery = () => {
  return useQuery({
    queryKey: queryKeys.user.me,
    queryFn: getUser,
    retry: false,
    staleTime: 5 * 60 * 1000, // Keep user data fresh for 5 mins
  });
};
