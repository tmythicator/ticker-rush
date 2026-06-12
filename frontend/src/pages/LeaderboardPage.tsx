import { Leaderboard } from '@/components/Leaderboard';
import { LadderDetails } from '@/components/Leaderboard/LadderDetails';
import { IconRefresh } from '@/components/icons/CustomIcons';
import { useActiveLadderQuery } from '@/hooks/useActiveLadderQuery';

export const LeaderboardPage = () => {
  const { data: ladder, isLoading } = useActiveLadderQuery();

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
        {ladder && <LadderDetails ladder={ladder} />}

        <div className="bg-background/20">
          <Leaderboard />
        </div>
      </div>
    </div>
  );
};
