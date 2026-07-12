import type { LeaderboardEntry } from '@/types';
import { LeaderboardHeader } from './LeaderboardHeader';
import { LeaderboardRow } from './LeaderboardRow';
import { LeaderboardEmptyState } from './LeaderboardEmptyState';

interface LeaderboardTableProps {
  entries: LeaderboardEntry[];
}

export const LeaderboardTable = ({ entries }: LeaderboardTableProps) => {
  return (
    <div className="overflow-hidden rounded-xl border border-border bg-card/50 shadow-sm">
      <table className="w-full text-left text-sm">
        <LeaderboardHeader />
        <tbody className="divide-y divide-border">
          {entries.map((entry) => (
            <LeaderboardRow key={entry.rank} entry={entry} />
          ))}
          {entries.length === 0 && <LeaderboardEmptyState />}
        </tbody>
      </table>
    </div>
  );
};
