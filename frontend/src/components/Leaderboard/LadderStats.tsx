import { IconCalendar, IconWallet } from '@/components/icons/CustomIcons';
import { formatLocalTime } from '@/lib/utils';

interface LadderStatsProps {
  endTime?: Date;
  initialBalance?: number;
}

export const LadderStats = ({ endTime, initialBalance }: LadderStatsProps) => {
  return (
    <div className="flex flex-wrap gap-4">
      <div className="flex items-center gap-4 bg-muted/30 p-4 rounded-xl border border-border">
        <div className="text-blue-500">
          <IconCalendar className="w-6 h-6" />
        </div>
        <div>
          <div className="text-[10px] uppercase font-black text-muted-foreground tracking-tighter mb-0.5">
            Competition Ends
          </div>
          <div className="text-base font-bold tabular-nums">
            {endTime ? formatLocalTime(endTime.getTime() / 1000) : 'N/A'}
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4 bg-muted/30 p-4 rounded-xl border border-border">
        <div className="text-emerald-500">
          <IconWallet className="w-6 h-6" />
        </div>
        <div>
          <div className="text-[10px] uppercase font-black text-muted-foreground tracking-tighter mb-0.5">
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
