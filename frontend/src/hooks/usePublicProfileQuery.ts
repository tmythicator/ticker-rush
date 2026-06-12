import { useQuery } from '@tanstack/react-query';
import { getPublicProfile } from '@/lib/api';
import { QUERY_KEY_PUBLIC_PROFILE } from '@/lib/queryKeys';

export const usePublicProfileQuery = (username?: string) => {
  return useQuery({
    queryKey: QUERY_KEY_PUBLIC_PROFILE(username || ''),
    queryFn: () => getPublicProfile({ username: username! }),
    enabled: !!username,
    retry: false,
  });
};
