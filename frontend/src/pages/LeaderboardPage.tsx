import { Leaderboard } from '@/components/Leaderboard';
import { LadderDetails } from '@/components/Leaderboard/LadderDetails';
import { IconRefresh } from '@/components/icons/CustomIcons';
import { useActiveLadderQuery } from '@/hooks/useActiveLadderQuery';

export const LeaderboardPage = () => {
  const { data: ladder, isLoading } = useActiveLadderQuery();

  if (isLoading) {
    return (
      <div className="container mx-auto flex h-64 items-center justify-center py-10">
        <IconRefresh className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-10">
      <div className="overflow-hidden rounded-2xl border border-border bg-card text-card-foreground shadow-sm">
        {ladder && <LadderDetails ladder={ladder} />}

        <div className="bg-background/20">
          <Leaderboard />
        </div>
      </div>
    </div>
  );
};
