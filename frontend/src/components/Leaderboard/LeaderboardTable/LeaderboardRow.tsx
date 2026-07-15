import { Link } from 'react-router-dom';
import { IconLock, IconMedal } from '@/components/icons/CustomIcons';
import type { LeaderboardEntry } from '@/types';
import styles from './LeaderboardTable.module.css';
import clsx from 'clsx';

interface LeaderboardRowProps {
  entry: LeaderboardEntry;
}

const MEDAL_CLASSES: Record<number, string> = {
  1: styles.medalGold,
  2: styles.medalSilver,
  3: styles.medalBronze,
};

export const LeaderboardRow = ({ entry }: LeaderboardRowProps) => {
  const medalClass = MEDAL_CLASSES[entry.rank];

  return (
    <tr className={styles.row}>
      <td className={styles.cellRank}>
        <div className={styles.rankContainer}>
          {medalClass ? <IconMedal className={clsx(styles.medalIcon, medalClass)} /> : entry.rank}
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

      <td className={clsx(styles.cell, styles.cellNetWorth)}>
        $
        {entry.score?.toLocaleString(undefined, {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2,
        })}
      </td>
    </tr>
  );
};
