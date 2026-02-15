import { IconMedal, IconRefresh, IconTrophy } from '@/components/icons/CustomIcons';
import { getLeaderboard } from '@/lib/api';
import { formatLocalTime } from '@/lib/utils';
import { useQuery } from '@tanstack/react-query';

export function Leaderboard() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['leaderboard'],
    queryFn: () => getLeaderboard(10, 0),
    refetchInterval: 60000, // Refresh every minute
  });

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
    <div className="w-full max-w-2xl mx-auto p-4 sm:p-6 space-y-6">
      <div className="flex items-center space-x-3 mb-6">
        <IconTrophy className="h-8 w-8 text-yellow-500" />
        <h2 className="text-2xl font-bold tracking-tight">Top Traders</h2>
      </div>

      <div className="rounded-lg border bg-card text-card-foreground shadow-sm overflow-hidden">
        <table className="w-full text-sm text-left">
          <thead className="bg-muted/50 text-muted-foreground uppercase text-xs font-semibold">
            <tr>
              <th className="px-6 py-3 w-16 text-center">#</th>
              <th className="px-6 py-3">Trader</th>
              <th className="px-6 py-3 text-right">Net Worth</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {data?.entries?.map((entry) => (
              <tr key={entry.user_id} className="hover:bg-muted/50 transition-colors">
                <td className="px-6 py-4 text-center font-medium flex justify-center items-center">
                  {entry.rank === 1 && <IconMedal className="h-5 w-5 text-yellow-500" />}
                  {entry.rank === 2 && <IconMedal className="h-5 w-5 text-gray-400" />}
                  {entry.rank === 3 && <IconMedal className="h-5 w-5 text-amber-700" />}
                  {entry.rank > 3 && entry.rank}
                </td>
                <td className="px-6 py-4 font-medium">
                  {entry.first_name} {entry.last_name}
                </td>
                <td className="px-6 py-4 text-right font-mono text-emerald-500 font-bold">
                  $
                  {entry.total_net_worth?.toLocaleString(undefined, {
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                  })}
                </td>
              </tr>
            ))}
            {(!data?.entries || data.entries.length === 0) && (
              <tr>
                <td colSpan={3} className="px-6 py-8 text-center text-muted-foreground">
                  No traders ranked yet. Be the first!
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <p className="text-xs text-center text-muted-foreground mt-4">
        Leaderboard updates every 10 minutes. Status refreshes automatically. Last update:{' '}
        {data ? formatLocalTime(data.last_update) : 'Never'}
      </p>
    </div>
  );
}
