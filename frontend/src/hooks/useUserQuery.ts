import { useQuery } from '@tanstack/react-query';
import { getUser } from '@/lib/api';
import { queryKeys } from '@/lib/queryKeys';
import { queryConfig } from '@/lib/queryConfig';

export const useUserQuery = () => {
  return useQuery({
    queryKey: queryKeys.user.me,
    queryFn: getUser,
    retry: false,
    ...queryConfig.user,
  });
};
