import { IconLock, IconMedal, IconRefresh, IconTrophy } from '@/components/icons/CustomIcons';
import { getLeaderboard } from '@/lib/api';
import { QUERY_KEY_LEADERBOARD } from '@/lib/queryKeys';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';

export function Leaderboard() {
  const { data, isLoading, error } = useQuery({
    queryKey: QUERY_KEY_LEADERBOARD,
    queryFn: () => getLeaderboard({ limit: 10, offset: 0 }),
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
    <div className="w-full mx-auto p-4 sm:p-8 space-y-6">
      <div className="flex items-center space-x-3 mb-6">
        <IconTrophy className="h-8 w-8 text-yellow-500" />
        <h2 className="text-2xl font-bold tracking-tight">Top Traders</h2>
      </div>

      <div className="rounded-xl border border-border bg-card/50 shadow-sm overflow-hidden">
        <table className="w-full text-sm text-left">
          <thead className="bg-muted text-muted-foreground uppercase text-[10px] font-black tracking-widest border-b border-border">
            <tr>
              <th className="px-8 py-4 w-24 text-center">Rank</th>
              <th className="px-6 py-4">Trader</th>
              <th className="px-8 py-4 text-right">Net Worth (USD)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {data?.entries?.map((entry) => (
              <tr key={entry.user?.id} className="hover:bg-muted/50 transition-colors">
                <td className="px-6 py-4 text-center font-bold">
                  <div className="flex justify-center items-center">
                    {entry.rank === 1 && <IconMedal className="h-5 w-5 text-yellow-500" />}
                    {entry.rank === 2 && <IconMedal className="h-5 w-5 text-slate-400" />}
                    {entry.rank === 3 && <IconMedal className="h-5 w-5 text-amber-700" />}
                    {entry.rank > 3 && entry.rank}
                  </div>
                </td>
                <td className="px-6 py-4 font-bold">
                  {entry.user?.is_public ? (
                    <Link
                      to={`/users/${entry.user?.username}`}
                      className="hover:text-primary transition-colors flex items-center gap-2"
                    >
                      {entry.user?.username}
                    </Link>
                  ) : (
                    <div className="flex items-center gap-2 text-muted-foreground/60 italic">
                      <span>Classified</span>
                      <IconLock className="w-3 h-3" />
                    </div>
                  )}
                </td>

                <td className="px-6 py-4 text-right font-bold tabular-nums text-emerald-500">
                  $
                  {entry.score?.toLocaleString(undefined, {
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                  })}
                </td>
              </tr>
            ))}
            {(!data?.entries || data.entries.length === 0) && (
              <tr>
                <td colSpan={3} className="px-6 py-8 text-center text-muted-foreground">
                  The arena is empty. Be the first!
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <p className="text-xs text-center text-muted-foreground mt-4 opacity-50">
        Leaderboard updates every 10 minutes. Status refreshes automatically.
      </p>
    </div>
  );
}
