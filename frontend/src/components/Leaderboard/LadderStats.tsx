import { IconCalendar, IconWallet } from '@/components/icons/CustomIcons';
import { formatLocalTime } from '@/lib/utils';

interface LadderStatsProps {
  endTime?: Date;
  initialBalance?: number;
}

export const LadderStats = ({ endTime, initialBalance }: LadderStatsProps) => {
  return (
    <div className="flex flex-wrap gap-4">
      <div className="flex items-center gap-4 rounded-xl border border-border bg-muted/30 p-4">
        <div className="text-blue-500">
          <IconCalendar className="h-6 w-6" />
        </div>
        <div>
          <div className="mb-0.5 text-[10px] font-black uppercase tracking-tighter text-muted-foreground">
            Competition Ends
          </div>
          <div className="text-base font-bold tabular-nums">
            {endTime ? formatLocalTime(endTime.getTime() / 1000) : 'N/A'}
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4 rounded-xl border border-border bg-muted/30 p-4">
        <div className="text-emerald-500">
          <IconWallet className="h-6 w-6" />
        </div>
        <div>
          <div className="mb-0.5 text-[10px] font-black uppercase tracking-tighter text-muted-foreground">
            Starting Capital
          </div>
          <div className="text-base font-bold tabular-nums">
            ${initialBalance?.toLocaleString() ?? '0'}
          </div>
        </div>
      </div>
    </div>
  );
};
