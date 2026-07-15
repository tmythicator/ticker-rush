import { LadderHeader } from './LadderHeader';
import { LadderStats } from './LadderStats';
import { LeaderBoardAssets } from './LeaderBoardAssets';
import type { Ladder } from '@/types';
import styles from './LadderDetails.module.css';

interface LadderDetailsProps {
  ladder: Ladder;
}

export const LadderDetails = ({ ladder }: LadderDetailsProps) => {
  return (
    <div className={styles.container}>
      <div className={styles.ladderMetaRow}>
        <LadderHeader name={ladder.name} type={ladder.type} />
        <LadderStats endTime={ladder.end_time} initialBalance={ladder.initial_balance} />
      </div>
      <LeaderBoardAssets assets={ladder.allowed_tickers} />
    </div>
  );
};
