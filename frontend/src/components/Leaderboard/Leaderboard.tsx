import { IconRefresh, IconTrophy } from '@/components/icons/CustomIcons';
import { useLeaderboardQuery } from '@/hooks/useLeaderboardQuery';
import { LeaderboardTable } from './LeaderboardTable';

export function Leaderboard() {
  const { data, isLoading, error } = useLeaderboardQuery();

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <IconRefresh className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error) {
    return <div className="text-center text-destructive p-4">Failed to load leaderboard.</div>;
  }

  return (
    <div className="w-full mx-auto p-4 sm:p-8 space-y-6">
      <div className="flex items-center space-x-3 mb-6">
        <IconTrophy className="h-8 w-8 text-yellow-500" />
        <h2 className="text-2xl font-bold tracking-tight">Top Traders</h2>
      </div>

      <LeaderboardTable entries={data?.entries ?? []} />

      <p className="text-xs text-center text-muted-foreground mt-4 opacity-50">
        Leaderboard updates regularly. Status refreshes automatically.
      </p>
    </div>
  );
}
