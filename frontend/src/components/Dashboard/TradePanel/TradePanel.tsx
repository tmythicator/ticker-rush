import { useState } from 'react';
import { useTrade } from '../../../hooks/useTrade';
import { TradeAction } from '../../../types';
import { TradeFooter } from './TradeFooter';
import { TradeOrderInput } from './TradeOrderInput';
import { TradePanelHeader } from './TradePanelHeader';
import { useAuth } from '../../../hooks/useAuth';

export interface TradePanelProps {
    symbol: string;
    currentPrice?: number;
    onTradeSuccess?: () => void;
}

export const TradePanel = ({ symbol, currentPrice = 0, onTradeSuccess }: TradePanelProps) => {
    const [quantity, setQuantity] = useState<string>('');
    const { user } = useAuth();
    const buyingPower = user?.balance || 0;

    const { executeTrade, isLoading, error } = useTrade({
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
