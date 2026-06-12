import { IconTrophy } from '@/components/icons/CustomIcons';
import { useJoinLadder } from '@/hooks/useJoinLadder';
import { JoinLadderActions } from './JoinLadderActions';
import { JoinLadderNotice } from './JoinLadderNotice';

export const JoinLadderButton = () => {
  const { isConfirming, setIsConfirming, isPending, handleJoin } = useJoinLadder();

  return (
    <div className="relative group overflow-hidden bg-card/30 backdrop-blur-md border border-primary/20 rounded-xl p-6 transition-all duration-300 hover:border-primary/50 hover:shadow-[0_0_20px_rgba(190,100,50,0.15)]">
      <div className="absolute top-0 right-0 -m-4 w-24 h-24 bg-primary/5 rounded-full blur-2xl group-hover:bg-primary/10 transition-colors duration-500" />

      <div className="flex flex-col md:flex-row items-center gap-6 relative z-10">
        <div className="flex-shrink-0 w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center text-primary group-hover:scale-110 transition-transform duration-300">
          <IconTrophy className="w-8 h-8" />
        </div>

        <div className="flex-grow text-center md:text-left">
          <h3 className="text-xl font-bold text-foreground mb-1">Join the Season</h3>
          <p className="text-muted-foreground text-sm max-w-md">
            Initialize your trading balance and compete for the top spot on the leaderboard.
          </p>
        </div>

        <div className="flex-shrink-0 w-full md:w-auto mt-4 md:mt-0">
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
