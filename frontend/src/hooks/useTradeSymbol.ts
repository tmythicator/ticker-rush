import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { TradeSymbol, isTradeSymbol } from '../types';

export const useTradeSymbol = () => {
    const [searchParams, setSearchParams] = useSearchParams();

    const initialSymbolParam = searchParams.get('symbol');
    const symbol: TradeSymbol = (initialSymbolParam && isTradeSymbol(initialSymbolParam))
        ? (initialSymbolParam as TradeSymbol)
        : TradeSymbol.AAPL;

    const setSymbol = (newSymbol: TradeSymbol) => {
        setSearchParams({ symbol: newSymbol });
    };

    useEffect(() => {
        if (!initialSymbolParam || !isTradeSymbol(initialSymbolParam)) {
            setSearchParams({ symbol: TradeSymbol.AAPL }, { replace: true });
        }
    }, [initialSymbolParam, setSearchParams]);

    return { symbol, setSymbol };
};
