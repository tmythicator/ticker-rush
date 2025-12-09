import { useState } from 'react';
import { useTrade } from '../../hooks/useTrade';
import { TradeAction } from '../../types';
import { TradeFooter } from './TradeFooter';
import { TradeOrderInput } from './TradeOrderInput';
import { TradePanelHeader } from './TradePanelHeader';

export interface TradePanelProps {
    userId: number;
    symbol: string;
    currentPrice?: number;
    buyingPower?: number;
    onTradeSuccess?: () => void;
}

export const TradePanel = ({ userId, symbol, currentPrice = 0, buyingPower = 0, onTradeSuccess }: TradePanelProps) => {
    const [quantity, setQuantity] = useState<string>('');

    const { executeTrade, isLoading, error } = useTrade({
        userId,
        symbol,
        onSuccess: () => {
            setQuantity('');
            if (onTradeSuccess) onTradeSuccess();
        }
    });

    const qty = parseFloat(quantity) || 0;
    const estCost = qty * currentPrice;

    const handleTrade = (action: TradeAction) => {
        executeTrade(action, qty);
    };

    return (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 flex flex-col h-full">
            <TradePanelHeader isLoading={isLoading} />
            <TradeOrderInput symbol={symbol} quantity={quantity} setQuantity={setQuantity} error={error} handleTrade={handleTrade} />
            <TradeFooter buyingPower={buyingPower} estCost={estCost} />
        </div>
    );
};
