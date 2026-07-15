import styles from './Leaderboard.module.css';

export const LeaderboardHeader = () => {
  return (
    <thead className={styles.thead}>
      <tr>
        <th className={`${styles.headerCell} ${styles.headerCellCenter}`}>Rank</th>
        <th className={styles.headerCell}>Trader</th>
        <th className={`${styles.headerCell} ${styles.headerCellRight}`}>Net Worth (USD)</th>
      </tr>
    </thead>
  );
};
