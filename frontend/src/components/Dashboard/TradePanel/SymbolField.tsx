import { SourceBadge } from '@/components/shared/SourceBadge';
import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import type { TickerSource } from '@/types';

interface SymbolFieldProps {
  symbol: string;
  source?: TickerSource;
}

export const SymbolField = ({ symbol, source }: SymbolFieldProps) => (
  <div>
    <Label className="block text-xs text-muted-foreground uppercase tracking-wider mb-2">
      Symbol
    </Label>
    <div className="w-full bg-muted border border-border rounded-lg px-3 py-3 flex items-center gap-3 opacity-70">
      {source && <SourceBadge source={source} />}
      <Input
        type="text"
        value={symbol}
        disabled
        variant="unstyled"
        size="unstyled"
        className="flex-1 font-mono text-sm font-bold text-muted-foreground cursor-default"
      />
    </div>
  </div>
);
