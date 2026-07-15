import styles from './LeaderboardTable.module.css';

export const LeaderboardEmptyState = () => {
  return (
    <tr>
      <td colSpan={3} className={styles.cellEmpty}>
        The arena is empty. Be the first!
      </td>
    </tr>
  );
};
