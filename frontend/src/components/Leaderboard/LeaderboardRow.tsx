import { Link } from 'react-router-dom';
import { IconLock, IconMedal } from '@/components/icons/CustomIcons';
import type { LeaderboardEntry } from '@/types';

interface LeaderboardRowProps {
  entry: LeaderboardEntry;
}

export const LeaderboardRow = ({ entry }: LeaderboardRowProps) => {
  return (
    <tr className="hover:bg-muted/50 transition-colors">
      <td className="px-6 py-4 text-center font-bold">
        <div className="flex justify-center items-center">
          {entry.rank === 1 && <IconMedal className="h-5 w-5 text-yellow-500" />}
          {entry.rank === 2 && <IconMedal className="h-5 w-5 text-slate-400" />}
          {entry.rank === 3 && <IconMedal className="h-5 w-5 text-amber-700" />}
          {entry.rank > 3 && entry.rank}
        </div>
      </td>
      <td className="px-6 py-4 font-bold">
        {entry.user?.is_public ? (
          <Link
            to={`/users/${entry.user?.username}`}
            className="hover:text-primary transition-colors flex items-center gap-2"
          >
            {entry.user?.username}
          </Link>
        ) : (
          <div className="flex items-center gap-2 text-muted-foreground/60 italic">
            <span>Classified</span>
            <IconLock className="w-3 h-3" />
          </div>
        )}
      </td>

      <td className="px-6 py-4 text-right font-bold tabular-nums text-emerald-500">
        $
        {entry.score?.toLocaleString(undefined, {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2,
        })}
      </td>
    </tr>
  );
};
