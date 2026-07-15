import { Link } from 'react-router-dom';
import { IconLock, IconMedal } from '@/components/icons/CustomIcons';
import type { LeaderboardEntry } from '@/types';
import styles from './Leaderboard.module.css';

interface LeaderboardRowProps {
  entry: LeaderboardEntry;
}

export const LeaderboardRow = ({ entry }: LeaderboardRowProps) => {
  return (
    <tr>
      <td className={`${styles.cell} ${styles.cellCenter}`}>
        <div className={styles.rankContainer}>
          {entry.rank === 1 && <IconMedal className={`${styles.medalIcon} ${styles.gold}`} />}
          {entry.rank === 2 && <IconMedal className={`${styles.medalIcon} ${styles.silver}`} />}
          {entry.rank === 3 && <IconMedal className={`${styles.medalIcon} ${styles.bronze}`} />}
          {entry.rank > 3 && entry.rank}
        </div>
      </td>
      <td className={styles.cell}>
        {entry.user?.is_public ? (
          <Link to={`/users/${entry.user?.username}`} className={styles.userLink}>
            {entry.user?.username}
          </Link>
        ) : (
          <div className={styles.classifiedContainer}>
            <span>Classified</span>
            <IconLock className={styles.lockIcon} />
          </div>
        )}
      </td>

      <td className={`${styles.cell} ${styles.cellRight}`}>
        $
        {entry.score?.toLocaleString(undefined, {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2,
        })}
      </td>
    </tr>
  );
};
