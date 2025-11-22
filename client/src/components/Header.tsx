import { Activity, Wallet } from 'lucide-react';

export const Header = () => {
    return (
        <header className="h-14 bg-white border-b border-slate-200 flex items-center px-6 justify-between sticky top-0 z-50">
            <div className="flex items-center gap-2">
                <div className="bg-blue-600 p-1.5 rounded-lg shadow-blue-100">
                    <Activity className="w-5 h-5 text-white" />
                </div>
                <h1 className="font-bold text-lg tracking-tight text-slate-900">Ticker Rush</h1>
            </div>

            <div className="flex items-center gap-6 text-sm font-medium">
                <div className="group flex items-center gap-2 text-slate-600 hover:text-blue-600 cursor-pointer transition-colors bg-white px-3 py-1.5 rounded-full border border-transparent hover:border-slate-200">
                    <Wallet className="w-4 h-4 text-slate-400 group-hover:text-blue-500" />
                    <span className="tabular-nums">$10,000.00</span>
                </div>
                <div className="w-8 h-8 bg-slate-100 rounded-full border border-slate-200 flex items-center justify-center text-slate-600 font-bold text-xs cursor-pointer hover:bg-slate-200 transition-colors">
                    JD
                </div>
            </div>
        </header>
    );
};