import { IconChevronDown } from '../../icons/CustomIcons';
import { isTradeSymbol, TradeSymbol, TradeSymbols } from '../../../types';

interface ChartSymbolPickerProps {
  symbol: string;
  onSymbolChange: (symbol: TradeSymbol) => void;
}

export const ChartSymbolPicker = ({ symbol, onSymbolChange }: ChartSymbolPickerProps) => {
  const handleSymbolChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newSymbol = e.target.value;
    if (isTradeSymbol(newSymbol)) {
      onSymbolChange(newSymbol);
    }
  };
  return (
    <div className="relative group/select">
      <select
        value={symbol}
        onChange={handleSymbolChange}
        className="appearance-none bg-transparent pl-3 pr-8 py-1.5 font-bold text-foreground text-lg tracking-tight focus:outline-none cursor-pointer hover:text-primary transition-colors"
      >
        {TradeSymbols.map((t) => (
          <option key={t} value={t} className="bg-popover text-popover-foreground">
            {t}
          </option>
        ))}
      </select>
      <IconChevronDown className="w-4 h-4 text-muted-foreground absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none group-hover/select:text-primary transition-colors" />
    </div>
  );
};
