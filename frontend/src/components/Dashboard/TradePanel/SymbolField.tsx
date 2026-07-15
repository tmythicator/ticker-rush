import { SourceBadge } from '@/components/shared/SourceBadge';
import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import type { TickerSource } from '@/types';
import styles from './TradePanel.module.css';

interface SymbolFieldProps {
  symbol: string;
  source?: TickerSource;
}

export const SymbolField = ({ symbol, source }: SymbolFieldProps) => (
  <div>
    <Label className={styles.label}>
      Symbol
    </Label>
    <div className={styles.inputWrapper}>
      {source && <SourceBadge source={source} />}
      <Input
        type="text"
        value={symbol}
        disabled
        variant="unstyled"
        size="unstyled"
        className={styles.inputDisabledText}
      />
    </div>
  </div>
);
