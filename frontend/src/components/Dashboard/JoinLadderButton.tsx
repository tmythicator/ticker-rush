import { IconLock, IconTrophy } from '@/components/icons/CustomIcons';
import { joinLadder } from '@/lib/api';
import { QUERY_KEY_USER } from '@/lib/queryKeys';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';

export const JoinLadderButton = () => {
  const queryClient = useQueryClient();
  const [isConfirming, setIsConfirming] = useState(false);

  const mutation = useMutation({
    mutationFn: joinLadder,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: QUERY_KEY_USER });
    },
  });

  const handleJoin = () => {
    mutation.mutate();
  };

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
          {!isConfirming ? (
            <button
              onClick={() => setIsConfirming(true)}
              className="w-full md:w-auto bg-primary text-primary-foreground font-bold py-3 px-8 rounded-full shadow-lg shadow-primary/20 hover:scale-[1.02] active:scale-[0.98] transition-all"
            >
              Get Started
            </button>
          ) : (
            <div className="flex flex-col gap-3">
              <button
                disabled={mutation.isPending}
                onClick={handleJoin}
                className="w-full md:w-auto bg-primary text-primary-foreground font-bold py-3 px-8 rounded-full shadow-lg shadow-primary/20 hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-50"
              >
                {mutation.isPending ? 'Joining...' : 'Confirm Entry'}
              </button>
              <button
                disabled={mutation.isPending}
                onClick={() => setIsConfirming(false)}
                className="text-xs text-muted-foreground hover:text-foreground transition-colors uppercase tracking-wider font-semibold"
              >
                Cancel
              </button>
            </div>
          )}
        </div>
      </div>

      <div className="mt-6 pt-6 border-t border-border flex items-start gap-4">
        <div className="mt-1">
          <IconLock className="w-4 h-4 text-secondary" />
        </div>
        <div className="text-xs space-y-2">
          <p className="text-muted-foreground italic">
            <span className="text-secondary font-bold">Important:</span> Once you join, your
            participation in this ladder cycle is permanent and cannot be undone. This ensures the
            integrity of the leaderboard and fair competition.
          </p>
          <div className="flex items-center gap-2 text-primary/80">
            <div className="w-1 h-1 bg-primary rounded-full" />
            <p>Privacy concern? You can always toggle your profile to Private in the settings.</p>
          </div>
        </div>
      </div>
    </div>
  );
};
