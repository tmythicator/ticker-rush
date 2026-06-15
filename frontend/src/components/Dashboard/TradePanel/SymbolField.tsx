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
    <Label className="mb-2 block text-xs uppercase tracking-wider text-muted-foreground">
      Symbol
    </Label>
    <div className="flex w-full items-center gap-3 rounded-lg border border-border bg-muted px-3 py-3 opacity-70">
      {source && <SourceBadge source={source} />}
      <Input
        type="text"
        value={symbol}
        disabled
        variant="unstyled"
        size="unstyled"
        className="flex-1 cursor-default font-mono text-sm font-bold text-muted-foreground"
      />
    </div>
  </div>
);
