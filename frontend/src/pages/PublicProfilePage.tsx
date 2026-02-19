import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';

import { PortfolioTable } from '@/components/PortfolioTable/PortfolioTable';
import { StatsGrid } from '@/components/Profile/StatsGrid';
import { IconLock, IconRefresh } from '@/components/icons/CustomIcons';
import { QUERY_KEY_PUBLIC_PROFILE } from '@/lib/queryKeys';
import { getPublicProfile } from '@/types';

export const PublicProfilePage = () => {
  const { username } = useParams<{ username: string }>();

  const {
    data: user,
    isLoading,
    error,
  } = useQuery({
    queryKey: QUERY_KEY_PUBLIC_PROFILE(username!),
    queryFn: () => getPublicProfile(username!),
    enabled: !!username,
    retry: false,
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <IconRefresh className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-4">
        <IconLock className="w-8 h-8 text-primary" />
        <h1 className="text-2xl font-bold text-foreground">Profile Unavailable</h1>
        <p className="text-muted-foreground">This profile is private or does not exist.</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen relative overflow-hidden bg-background">
      <div className="fixed -top-[20%] -right-[10%] w-[800px] h-[800px] bg-primary/5 rounded-full blur-3xl pointer-events-none" />
      <div className="fixed top-[10%] left-[5%] w-[600px] h-[600px] bg-purple-500/5 rounded-full blur-3xl pointer-events-none" />

      <div className="container max-w-7xl mx-auto py-12 px-6 space-y-12 relative z-10">
        <div className="flex flex-col gap-2">
          <h1 className="text-3xl font-bold text-foreground bg-clip-text text-transparent bg-gradient-to-r from-primary to-purple-600 w-fit">
            {user.first_name} {user.last_name}
          </h1>
          <div className="flex flex-col gap-1 text-muted-foreground">
            <span className="font-mono">@{user.username}</span>
            {user.website && (
              <div className="flex items-center gap-2 mt-1">
                <span className="text-sm font-semibold text-foreground/80">Website:</span>
                <a
                  href={user.website}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary hover:underline hover:text-primary/90 transition-colors"
                >
                  {user.website}
                </a>
              </div>
            )}
          </div>
        </div>

        <StatsGrid {...user} />

        <div className="space-y-4">
          <h2 className="text-xl font-semibold text-foreground">Portfolio</h2>
          <div className="glass-panel rounded-xl overflow-hidden border border-border/50">
            <PortfolioTable portfolio={user.portfolio || {}} isReadOnly />
          </div>
        </div>
      </div>
    </div>
  );
};
