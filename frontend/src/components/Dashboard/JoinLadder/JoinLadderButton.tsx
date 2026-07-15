import { IconTrophy } from '@/components/icons/CustomIcons';
import { useJoinLadder } from '@/hooks/useJoinLadder';
import { JoinLadderActions } from './JoinLadderActions';
import { JoinLadderNotice } from './JoinLadderNotice';
import styles from './JoinLadderButton.module.css';

export const JoinLadderButton = () => {
  const { isConfirming, setIsConfirming, isPending, handleJoin } = useJoinLadder();

  return (
    <div className={styles.card}>
      <div className={styles.glow} />

      <div className={styles.joinCardBody}>
        <div className={styles.joinIcon}>
          <IconTrophy />
        </div>

        <div className={styles.textGroup}>
          <h3 className={styles.title}>Join the Season</h3>
          <p className={styles.description}>
            Initialize your trading balance and compete for the top spot on the leaderboard.
          </p>
        </div>

        <div className={styles.actionWrapper}>
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
