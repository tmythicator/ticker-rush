import { useQuery } from '@tanstack/react-query';
import { getPublicProfile } from '@/lib/api';
import { queryKeys } from '@/lib/queryKeys';

export const usePublicProfileQuery = (username?: string) => {
  return useQuery({
    queryKey: queryKeys.user.publicProfile(username || ''),
    queryFn: () => getPublicProfile({ username: username! }),
    enabled: !!username,
    retry: false,
  });
};
