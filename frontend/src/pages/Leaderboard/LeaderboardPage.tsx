import { Leaderboard } from '@/components/Leaderboard';
import { LadderDetails } from '@/components/Leaderboard/LadderDetails';
import { IconRefresh } from '@/components/icons/CustomIcons';
import { Card } from '@/components/shared/Card';
import { useActiveLadderQuery } from '@/hooks/useActiveLadderQuery';
import styles from './LeaderboardPage.module.css';

export const LeaderboardPage = () => {
  const { data: ladder, isLoading } = useActiveLadderQuery();

  if (isLoading) {
    return (
      <div className={styles.loaderWrapper}>
        <IconRefresh className={styles.spinner} />
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <Card className={styles.card}>
        {ladder && <LadderDetails ladder={ladder} />}

        <div className={styles.leaderboardSection}>
          <Leaderboard />
        </div>
      </Card>
    </div>
  );
};
