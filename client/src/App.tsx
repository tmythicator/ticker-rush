import { DollarSign, Briefcase, Percent } from 'lucide-react';
import { Header } from './components/Header';
import { MarketChart } from './components/MarketChart';
import { StatCard } from './components/StatCard';
import { TradePanel } from './components/TradePanel';

function App() {
    // TODO: replace mock data from the backend
    const stats = [
        { label: "Daily P&L", value: "+$1,240.50", trend: "+12.5%", icon: DollarSign },
        { label: "Open Positions", value: "4", icon: Briefcase },
        { label: "Win Rate", value: "68%", icon: Percent },
    ];

    return (
        <div className="min-h-screen bg-slate-50 text-slate-900 flex flex-col font-sans">
            <Header />

            <main className="flex-1 max-w-[1800px] w-full mx-auto p-4 lg:p-6 grid grid-cols-1 lg:grid-cols-12 gap-6">
                {/* Left Zone: Chart & Analytics (9/12) */}
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

                {/* Right Zone: Action Panel (3/12) */}
                <div className="lg:col-span-3 flex flex-col gap-4 h-full">
                    <TradePanel />
                </div>
            </main>
        </div>
    );
}

export default App;