import {
  LadderHeader,
  LadderStats,
  Leaderboard,
  LeaderBoardAssets,
} from '@/components/Leaderboard';
import { IconRefresh } from '@/components/icons/CustomIcons';
import { getActiveLadder } from '@/lib/api';
import { useQuery } from '@tanstack/react-query';

export const LeaderboardPage = () => {
  const { data: ladder, isLoading } = useQuery({
    queryKey: ['ladder', 'active'],
    queryFn: getActiveLadder,
  });

  if (isLoading) {
    return (
      <div className="container mx-auto py-10 flex justify-center items-center h-64">
        <IconRefresh className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="container mx-auto py-10 px-4">
      <div className="bg-card rounded-2xl border border-border shadow-sm overflow-hidden text-card-foreground">
        {ladder && (
          <div className="p-8 border-b border-border">
            <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-8 mb-8">
              <LadderHeader name={ladder.name} type={ladder.type} />
              <LadderStats endTime={ladder.end_time} initialBalance={ladder.initial_balance} />
            </div>
            <LeaderBoardAssets assets={ladder.allowed_tickers} />
          </div>
        )}

        <div className="bg-background/20">
          <Leaderboard />
        </div>
      </div>
    </div>
  );
};
