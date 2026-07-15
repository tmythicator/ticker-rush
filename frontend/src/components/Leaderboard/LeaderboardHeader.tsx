import styles from './Leaderboard.module.css';

export const LeaderboardHeader = () => {
  return (
    <thead className={styles.thead}>
      <tr>
        <th className={styles.headerCell} data-align="center">Rank</th>
        <th className={styles.headerCell}>Trader</th>
        <th className={styles.headerCell} data-align="right">Net Worth (USD)</th>
      </tr>
    </thead>
  );
};
