import { type TradeSymbol } from '@/types';
import { ChartSymbolIndicator } from './ChartSymbolIndicator';
import { ChartSymbolPicker } from './ChartSymbolPicker';

interface ChartBodyProps {
  symbol: TradeSymbol;
  price: number | undefined;
  isClosed: boolean | undefined;
  priceColor: string;
  isLoading: boolean;
  isError: boolean;

  onSymbolChange: (s: TradeSymbol) => void;
}

export const ChartBody = ({
  symbol,
  onSymbolChange,
  price,
  isClosed,
  priceColor,
  isLoading,
  isError,
}: ChartBodyProps) => {
  return (
    <div className="absolute top-4 left-4 z-20 flex flex-col gap-2">
      <div className="bg-background/90 backdrop-blur-md p-1.5 rounded-xl shadow-lg border border-border/60 flex items-center gap-1 transition-all hover:shadow-xl hover:scale-[1.02]">
        <ChartSymbolPicker symbol={symbol} onSymbolChange={onSymbolChange} />
        <div className="h-6 w-px bg-border mx-1"></div>
        <ChartSymbolIndicator
          price={price}
          isClosed={isClosed}
          priceColor={priceColor}
          isLoading={isLoading}
          isError={isError}
        />
      </div>
    </div>
  );
};
