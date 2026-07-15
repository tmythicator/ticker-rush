import { IconRefresh, IconTrophy } from '@/components/icons/CustomIcons';
import { useLeaderboardQuery } from '@/hooks/useLeaderboardQuery';
import { LeaderboardTable } from './LeaderboardTable/LeaderboardTable';
import styles from './Leaderboard.module.css';

export function Leaderboard() {
  const { data, isLoading, error } = useLeaderboardQuery();

  if (isLoading) {
    return (
      <div className={styles.loaderWrapper}>
        <IconRefresh className={styles.spinner} />
      </div>
    );
  }

  if (error) {
    return <div className={styles.error}>Failed to load leaderboard.</div>;
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <IconTrophy className={styles.trophyIcon} />
        <h2 className={styles.title}>Top Traders</h2>
      </div>

      <LeaderboardTable entries={data?.entries ?? []} />

      <p className={styles.footer}>
        Leaderboard updates regularly. Status refreshes automatically.
      </p>
    </div>
  );
}
