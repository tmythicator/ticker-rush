import { SourceBadge } from '@/components/shared/SourceBadge';
import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import type { TickerSource } from '@/types';
import styles from './TradePanel.module.css';
import { useId } from 'react';

interface SymbolFieldProps {
  symbol: string;
  source?: TickerSource;
}

export const SymbolField = ({ symbol, source }: SymbolFieldProps) => {
  const inputId = useId();

  return (
    <div>
      <div className={styles.labelRow}>
        <Label htmlFor={inputId} className={styles.label}>
          Symbol
        </Label>
      </div>
      <div className={styles.inputWrapper}>
        {source && <SourceBadge source={source} />}
        <Input
          id={inputId}
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
};
