import { LadderHeader } from './LadderHeader';
import { LadderStats } from './LadderStats';
import { LeaderBoardAssets } from './LeaderBoardAssets';
import type { Ladder } from '@/types';

interface LadderDetailsProps {
  ladder: Ladder;
}

export const LadderDetails = ({ ladder }: LadderDetailsProps) => {
  return (
    <div className="p-8 border-b border-border">
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-8 mb-8">
        <LadderHeader name={ladder.name} type={ladder.type} />
        <LadderStats endTime={ladder.end_time} initialBalance={ladder.initial_balance} />
      </div>
      <LeaderBoardAssets assets={ladder.allowed_tickers} />
    </div>
  );
};
