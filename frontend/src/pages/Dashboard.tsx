import { useState, useEffect } from 'react';
import { DollarSign, Briefcase, Percent } from 'lucide-react';
import { MarketChart } from '../components/MarketChart';
import { StatCard } from '../components/StatCard';
import { TradePanel } from '../components/TradePanel';
import { getUser, fetchQuote, type User } from '../lib/api';

const USER_ID = 1; // Hardcoded for dev
const SYMBOL = 'AAPL';

export const Dashboard = () => {
    const [user, setUser] = useState<User | null>(null);
    const [price, setPrice] = useState<number | undefined>(undefined);

    const loadData = async () => {
        try {
            const [userData, quote] = await Promise.all([
                getUser(USER_ID),
                fetchQuote(SYMBOL)
            ]);
            setUser(userData);
            setPrice(quote.price);
        } catch (e) {
            console.error('Failed to load dashboard data', e);
        }
    };

    useEffect(() => {
        loadData();
        const interval = setInterval(loadData, 5000); // Poll every 5s
        return () => clearInterval(interval);
    }, []);

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
                    <MarketChart />
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    {stats.map((stat, i) => (
                        <StatCard key={i} {...stat} />
                    ))}
                </div>
            </div>
            <div className="lg:col-span-3 flex flex-col gap-4 h-full">
                <TradePanel
                    userId={USER_ID}
                    symbol={SYMBOL}
                    currentPrice={price}
                    buyingPower={user?.balance}
                    onTradeSuccess={loadData}
                />
            </div>
        </div>
    );
};