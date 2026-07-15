import type { LeaderboardEntry } from '@/types';
import { LeaderboardHeader } from './LeaderboardHeader';
import { LeaderboardRow } from './LeaderboardRow';
import { LeaderboardEmptyState } from './LeaderboardEmptyState';
import styles from './Leaderboard.module.css';

interface LeaderboardTableProps {
  entries: LeaderboardEntry[];
}

export const LeaderboardTable = ({ entries }: LeaderboardTableProps) => {
  return (
    <div className={styles.tableWrapper}>
      <table className={styles.table}>
        <LeaderboardHeader />
        <tbody className={styles.tbody}>
          {entries.map((entry) => (
            <LeaderboardRow key={entry.rank} entry={entry} />
          ))}
          {entries.length === 0 && <LeaderboardEmptyState />}
        </tbody>
      </table>
    </div>
  );
};
