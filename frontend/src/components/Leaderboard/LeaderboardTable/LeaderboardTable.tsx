import type { LeaderboardEntry } from '@/types';
import { LeaderboardRow } from './LeaderboardRow';
import { LeaderboardEmptyState } from './LeaderboardEmptyState';
import styles from './LeaderboardTable.module.css';

interface LeaderboardTableProps {
  entries: LeaderboardEntry[];
}

export const LeaderboardTable = ({ entries }: LeaderboardTableProps) => {
  return (
    <div className={styles.wrapper}>
      <table className={styles.root}>
        <thead className={styles.head}>
          <tr>
            <th className={styles.cellRank}>Rank</th>
            <th className={styles.headerCell}>Trader</th>
            <th className={styles.headerCell}>Net Worth (USD)</th>
          </tr>
        </thead>
        <tbody>
          {entries.map((entry) => (
            <LeaderboardRow key={entry.rank} entry={entry} />
          ))}
          {entries.length === 0 && <LeaderboardEmptyState />}
        </tbody>
      </table>
    </div>
  );
};
