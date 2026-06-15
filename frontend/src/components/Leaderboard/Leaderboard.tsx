import { IconRefresh, IconTrophy } from '@/components/icons/CustomIcons';
import { useLeaderboardQuery } from '@/hooks/useLeaderboardQuery';
import { LeaderboardTable } from './LeaderboardTable';

export function Leaderboard() {
  const { data, isLoading, error } = useLeaderboardQuery();

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <IconRefresh className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error) {
    return <div className="p-4 text-center text-destructive">Failed to load leaderboard.</div>;
  }

  return (
    <div className="mx-auto w-full space-y-6 p-4 sm:p-8">
      <div className="mb-6 flex items-center space-x-3">
        <IconTrophy className="h-8 w-8 text-yellow-500" />
        <h2 className="text-2xl font-bold tracking-tight">Top Traders</h2>
      </div>

      <LeaderboardTable entries={data?.entries ?? []} />

      <p className="mt-4 text-center text-xs text-muted-foreground opacity-50">
        Leaderboard updates regularly. Status refreshes automatically.
      </p>
    </div>
  );
}
