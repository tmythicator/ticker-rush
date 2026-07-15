import { LadderHeader } from './LadderHeader';
import { LadderStats } from './LadderStats';
import { LadderAssets } from './LadderAssets';
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
      <LadderAssets assets={ladder.allowed_tickers} />
    </div>
  );
};
