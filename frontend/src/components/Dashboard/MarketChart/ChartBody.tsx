import { TradeSymbol } from '../../../types';
import { ChartSymbolIndicator } from './ChartSymbolIndicator';
import { ChartSymbolPicker } from './ChartSymbolPicker';

interface ChartBodyProps {
    symbol: TradeSymbol;
    price: number | undefined;
    priceColor: string;
    isLoading: boolean;
    isError: boolean;

    onSymbolChange: (s: TradeSymbol) => void;
}

export const ChartBody = ({ symbol, onSymbolChange, price, priceColor, isLoading, isError }: ChartBodyProps) => {

    return (
        <div className="absolute top-4 left-4 z-20 flex flex-col gap-2">
            <div className="bg-white/90 backdrop-blur-md p-1.5 rounded-xl shadow-lg border border-slate-200/60 flex items-center gap-1 transition-all hover:shadow-xl hover:scale-[1.02]">
                <ChartSymbolPicker symbol={symbol} onSymbolChange={onSymbolChange} />
                <div className="h-6 w-px bg-slate-200 mx-1"></div>
                <ChartSymbolIndicator price={price} priceColor={priceColor} isLoading={isLoading} isError={isError} />
            </div>
        </div>
    );
};
