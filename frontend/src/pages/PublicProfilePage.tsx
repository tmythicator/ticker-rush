import { useParams } from 'react-router-dom';
import { PortfolioTable } from '@/components/PortfolioTable';
import { StatsGrid } from '@/components/Profile/StatsGrid';
import { IconLock, IconRefresh } from '@/components/icons/CustomIcons';
import { usePublicProfileQuery } from '@/hooks/usePublicProfileQuery';

export const PublicProfilePage = () => {
  const { username } = useParams<{ username: string }>();

  const { data: user, isLoading, error } = usePublicProfileQuery(username);

  if (isLoading) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <IconRefresh className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error || !user) {
    return (
      <div
        data-testid="profile-unavailable"
        className="flex min-h-[60vh] flex-col items-center justify-center gap-4"
      >
        <IconLock className="h-8 w-8 text-primary" />
        <h1 className="text-2xl font-bold text-foreground">Profile Unavailable</h1>
        <p className="text-muted-foreground">This profile is private or does not exist.</p>
      </div>
    );
  }

  return (
    <div className="relative min-h-screen overflow-hidden bg-background">
      <div className="pointer-events-none fixed -right-[10%] -top-[20%] h-[800px] w-[800px] rounded-full bg-primary/5 blur-3xl" />
      <div className="pointer-events-none fixed left-[5%] top-[10%] h-[600px] w-[600px] rounded-full bg-purple-500/5 blur-3xl" />

      <div className="container relative z-10 mx-auto max-w-7xl space-y-12 px-6 py-12">
        <div className="flex flex-col gap-2">
          <h1
            data-testid="profile-name"
            className="w-fit bg-gradient-to-r from-primary to-purple-600 bg-clip-text text-3xl font-bold text-foreground text-transparent"
          >
            {user.first_name} {user.last_name}
          </h1>
          <div className="flex flex-col gap-1 text-muted-foreground">
            <span data-testid="profile-username" className="font-mono">
              @{user.username}
            </span>
            {user.website && (
              <div className="mt-1 flex items-center gap-2">
                <span className="text-sm font-semibold text-foreground/80">Website:</span>
                <a
                  href={user.website}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary transition-colors hover:text-primary/90 hover:underline"
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
          <div className="glass-panel overflow-hidden rounded-xl border border-border/50">
            <PortfolioTable items={Object.values(user.portfolio || {})} isReadOnly />
          </div>
        </div>
      </div>
    </div>
  );
};
