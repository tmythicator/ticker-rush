import { useRef } from 'react';
import { type Quote } from '../../../lib/api';
import { ChartBody } from './ChartBody';
import { TradeSymbol } from '../../../types';
import { useChart } from '../../../hooks/useChart';
import { usePriceColor } from '../../../hooks/usePriceColor';

interface MarketChartProps {
    symbol: TradeSymbol;
    onSymbolChange: (s: TradeSymbol) => void;
    quote?: Quote;
    isLoading: boolean;
    isError: boolean;
}

export const MarketChart = ({ symbol, onSymbolChange, quote, isLoading, isError }: MarketChartProps) => {
    const chartContainerRef = useRef<HTMLDivElement>(null);

    useChart({ chartContainerRef, quote, symbol });

    const priceColor = usePriceColor(quote?.price);

    return (
        <div className="w-full h-full relative group">
            <ChartBody
                symbol={symbol}
                onSymbolChange={onSymbolChange}
                price={quote?.price}
                priceColor={priceColor}
                isLoading={isLoading}
                isError={isError}
            />
            <div className="w-full h-full relative group">
                <div ref={chartContainerRef} className="w-full h-[500px]" />
            </div>
        </div>
    );
};
