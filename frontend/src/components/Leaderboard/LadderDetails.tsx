import { LadderHeader } from './LadderHeader';
import { LadderStats } from './LadderStats';
import { LeaderBoardAssets } from './LeaderBoardAssets';
import type { Ladder } from '@/types';

interface LadderDetailsProps {
  ladder: Ladder;
}

export const LadderDetails = ({ ladder }: LadderDetailsProps) => {
  return (
    <div className="border-b border-border p-8">
      <div className="mb-8 flex flex-col justify-between gap-8 lg:flex-row lg:items-center">
        <LadderHeader name={ladder.name} type={ladder.type} />
        <LadderStats endTime={ladder.end_time} initialBalance={ladder.initial_balance} />
      </div>
      <LeaderBoardAssets assets={ladder.allowed_tickers} />
    </div>
  );
};
