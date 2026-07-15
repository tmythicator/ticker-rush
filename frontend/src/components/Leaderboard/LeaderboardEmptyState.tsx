import styles from './Leaderboard.module.css';

export const LeaderboardEmptyState = () => {
  return (
    <tr>
      <td colSpan={3} className={styles.emptyCell}>
        The arena is empty. Be the first!
      </td>
    </tr>
  );
};
