import { Button } from '@/components/shared/Button';
import styles from './JoinLadderActions.module.css';

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
      <div className={styles.container}>
        <Button onClick={() => setIsConfirming(true)} size="lg">
          Get Started
        </Button>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <Button disabled={isPending} onClick={onJoin} size="lg">
        {isPending ? 'Joining...' : 'Confirm Entry'}
      </Button>
      <Button
        disabled={isPending}
        onClick={() => setIsConfirming(false)}
        variant="ghost"
        className={styles.cancelButton}
      >
        Cancel
      </Button>
    </div>
  );
};
