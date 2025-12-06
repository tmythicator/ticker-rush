import { useState } from 'react';
import { DollarSign, Briefcase, Percent } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { MarketChart } from '../components/MarketChart';
import { StatCard } from '../components/StatCard';
import { TradePanel } from '../components/TradePanel';
import { getUser, fetchQuote } from '../lib/api';
import { TradeSymbol } from '../types';
const TEST_USER_ID = 1; // Hardcoded for dev

export const Dashboard = () => {
    const [symbol, setSymbol] = useState<TradeSymbol>(TradeSymbol.AAPL);

    const { data: user, refetch: refetchUser } = useQuery({
        queryKey: ['user', TEST_USER_ID],
        queryFn: () => getUser(TEST_USER_ID),
        refetchInterval: 1000,
    });

    const { data: quote, isLoading: isQuoteLoading, isError: isQuoteError } = useQuery({
        queryKey: ['quote', symbol],
        queryFn: () => fetchQuote(symbol),
        refetchInterval: 3000,
    });

    // TODO: replace mock stats with calculated stats from user portfolio
    const stats = [
        { label: "Daily P&L", value: "+$1,240.50", trend: "+12.5%", icon: DollarSign },
        { label: "Open Positions", value: Object.keys(user?.portfolio || {}).length.toString(), icon: Briefcase },
        { label: "Win Rate", value: "68%", icon: Percent },
    ];

    return (
        <div className="max-w-[1800px] w-full mx-auto p-4 lg:p-6 grid grid-cols-1 lg:grid-cols-12 gap-6">
            <div className="lg:col-span-9 flex flex-col gap-6">
                <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-1 overflow-hidden h-[500px] relative">
                    <MarketChart
                        symbol={symbol}
                        onSymbolChange={setSymbol}
                        quote={quote}
                        isLoading={isQuoteLoading}
                        isError={isQuoteError}
                    />
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    {stats.map((stat, i) => (
                        <StatCard key={i} {...stat} />
                    ))}
                </div>
            </div>
            <div className="lg:col-span-3 flex flex-col gap-4 h-full">
                <TradePanel
                    userId={TEST_USER_ID}
                    symbol={symbol}
                    currentPrice={quote?.price}
                    buyingPower={user?.balance}
                    onTradeSuccess={() => refetchUser()}
                />
            </div>
        </div>
    );
};