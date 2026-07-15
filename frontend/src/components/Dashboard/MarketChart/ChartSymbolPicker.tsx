import { IconChevronDown } from '@/components/icons/CustomIcons';
import { SourceBadge } from '@/components/shared/SourceBadge';
import { useTickers } from '@/hooks/useTickers';
import { isTradeSymbol, type TickerSource, type TradeSymbol } from '@/types';
import styles from './MarketChart.module.css';

interface ChartSymbolPickerProps {
  symbol: TradeSymbol | null;
  onSymbolChange: (symbol: TradeSymbol) => void;
}

export const ChartSymbolPicker = ({ symbol, onSymbolChange }: ChartSymbolPickerProps) => {
  const { data: config } = useTickers();
  const tickers = config || [];

  const handleSymbolChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newSymbol = e.target.value;
    if (isTradeSymbol(newSymbol, tickers)) {
      onSymbolChange(newSymbol);
    }
  };

  const tickerInfo = tickers.find((t) => t.symbol === symbol);
  const source = (tickerInfo?.source ?? 'Finnhub') as TickerSource;

  return (
    <div className={styles.pickerWrapper}>
      <div className={styles.pickerContainer}>
        {symbol && <SourceBadge source={source} />}
        <select
          value={symbol || ''}
          onChange={handleSymbolChange}
          disabled={tickers.length === 0}
          className={styles.select}
        >
          {tickers.length === 0 ? (
            <option value="" className={styles.option}>No assets available</option>
          ) : (
            tickers.map((t) => (
              <option
                key={t.symbol}
                value={t.symbol}
                className={styles.option}
              >
                {t.symbol.toUpperCase()}
              </option>
            ))
          )}
        </select>
        <IconChevronDown className={styles.chevron} />
      </div>
    </div>
  );
};
