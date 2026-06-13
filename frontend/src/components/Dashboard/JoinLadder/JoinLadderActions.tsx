import { Button } from '@/components/shared/Button';

interface JoinLadderActionsProps {
  isConfirming: boolean;
  setIsConfirming: (confirming: boolean) => void;
  isPending: boolean;
  onJoin: () => void;
}

export const JoinLadderActions = ({
  isConfirming,
  setIsConfirming,
  isPending,
  onJoin,
}: JoinLadderActionsProps) => {
  if (!isConfirming) {
    return (
      <Button onClick={() => setIsConfirming(true)} size="lg" className="w-full md:w-auto">
        Get Started
      </Button>
    );
  }

  return (
    <div className="flex flex-col gap-3">
      <Button disabled={isPending} onClick={onJoin} size="lg" className="w-full md:w-auto">
        {isPending ? 'Joining...' : 'Confirm Entry'}
      </Button>
      <Button
        disabled={isPending}
        onClick={() => setIsConfirming(false)}
        variant="ghost"
        className="text-xs text-muted-foreground hover:text-foreground uppercase tracking-wider font-semibold h-auto py-1.5"
      >
        Cancel
      </Button>
    </div>
  );
};
