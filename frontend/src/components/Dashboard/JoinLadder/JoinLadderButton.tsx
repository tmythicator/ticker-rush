import { IconTrophy } from '@/components/icons/CustomIcons';
import { useJoinLadder } from '@/hooks/useJoinLadder';
import { JoinLadderActions } from './JoinLadderActions';
import { JoinLadderNotice } from './JoinLadderNotice';

export const JoinLadderButton = () => {
  const { isConfirming, setIsConfirming, isPending, handleJoin } = useJoinLadder();

  return (
    <div className="group relative overflow-hidden rounded-xl border border-primary/20 bg-card/30 p-6 backdrop-blur-md transition-all duration-300 hover:border-primary/50 hover:shadow-[0_0_20px_rgba(190,100,50,0.15)]">
      <div className="absolute right-0 top-0 -m-4 h-24 w-24 rounded-full bg-primary/5 blur-2xl transition-colors duration-500 group-hover:bg-primary/10" />

      <div className="relative z-10 flex flex-col items-center gap-6 md:flex-row">
        <div className="flex h-16 w-16 flex-shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary transition-transform duration-300 group-hover:scale-110">
          <IconTrophy className="h-8 w-8" />
        </div>

        <div className="flex-grow text-center md:text-left">
          <h3 className="mb-1 text-xl font-bold text-foreground">Join the Season</h3>
          <p className="max-w-md text-sm text-muted-foreground">
            Initialize your trading balance and compete for the top spot on the leaderboard.
          </p>
        </div>

        <div className="mt-4 w-full flex-shrink-0 md:mt-0 md:w-auto">
          <JoinLadderActions
            isConfirming={isConfirming}
            setIsConfirming={setIsConfirming}
            isPending={isPending}
            onJoin={handleJoin}
          />
        </div>
      </div>

      <JoinLadderNotice />
    </div>
  );
};
