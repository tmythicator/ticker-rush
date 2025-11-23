import { History, ArrowUpRight, ArrowDownRight, Clock } from 'lucide-react';

// TODO: get real data from backend
const history = [
    { id: 1, symbol: "AAPL", type: "BUY", qty: 10, price: 150.50, date: "2023-10-24 10:30", pnl: null },
    { id: 2, symbol: "TSLA", type: "SELL", qty: 5, price: 240.20, date: "2023-10-23 14:15", pnl: "+$120.00" },
    { id: 3, symbol: "NVDA", type: "BUY", qty: 2, price: 420.00, date: "2023-10-22 09:00", pnl: null },
];

export const Profile = () => {
    return (
        <div className="max-w-4xl w-full mx-auto p-4 lg:p-6">
            <div className="mb-8">
                <h2 className="text-2xl font-bold text-slate-900">Trading Journal</h2>
                <p className="text-slate-500">Track your performance and history</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                <div className="bg-slate-900 text-white p-6 rounded-xl shadow-lg">
                    <div className="text-slate-400 text-sm font-bold uppercase mb-1">Net Worth</div>
                    <div className="text-3xl font-bold">$11,240.50</div>
                    <div className="text-green-400 text-sm mt-2 flex items-center gap-1">
                        <ArrowUpRight className="w-3 h-3" /> +12.4% this month
                    </div>
                </div>
            </div>

            <div className="bg-white rounded-xl shadow-sm border border-slate-200">
                <div className="px-6 py-4 border-b border-slate-100 flex items-center gap-2">
                    <History className="w-4 h-4 text-slate-400" />
                    <h3 className="font-bold text-slate-800">Recent Transactions</h3>
                </div>
                <div className="divide-y divide-slate-50">
                    {history.map((tx) => (
                        <div key={tx.id} className="px-6 py-4 flex items-center justify-between hover:bg-slate-50 transition-colors">
                            <div className="flex items-center gap-4">
                                <div className={`p-2 rounded-lg ${tx.type === 'BUY' ? 'bg-green-50' : 'bg-red-50'}`}>
                                    {tx.type === 'BUY'
                                        ? <ArrowDownRight className="w-5 h-5 text-green-600" />
                                        : <ArrowUpRight className="w-5 h-5 text-red-600" />
                                    }
                                </div>
                                <div>
                                    <div className="font-bold text-slate-900">{tx.symbol} <span className="text-slate-400 font-normal ml-1">{tx.type}</span></div>
                                    <div className="text-xs text-slate-500 flex items-center gap-1">
                                        <Clock className="w-3 h-3" /> {tx.date}
                                    </div>
                                </div>
                            </div>
                            <div className="text-right">
                                <div className="font-mono font-bold text-slate-900">{tx.qty} @ ${tx.price}</div>
                                {tx.pnl && (
                                    <div className="text-xs font-bold text-green-600">{tx.pnl}</div>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};