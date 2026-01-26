import { Briefcase, Wallet, TrendingUp } from 'lucide-react';

import { type User } from '../../lib/api';
import { calculateInvestedCapital } from '../../lib/utils';

interface StatsGridProps {
  user: User;
}

export const StatsGrid = ({ user }: StatsGridProps) => {
  const portfolioItems = Object.values(user.portfolio ?? {});
  const investedCapital = calculateInvestedCapital(user.portfolio);
  const totalNetWorth = user.balance + investedCapital;
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div className="bg-slate-900 text-white p-6 rounded-2xl shadow-xl reltive overflow-hidden group">
        <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <Wallet className="w-24 h-24" />
        </div>
        <div className="text-slate-400 text-sm font-bold uppercase tracking-wider mb-2">
          Total Net Worth
        </div>
        <div className="text-4xl font-bold font-mono tracking-tight">
          ${totalNetWorth.toFixed(2)}
        </div>
        <div className="mt-4 flex items-center gap-2 text-sm text-slate-300">
          <span className="bg-slate-800 px-2 py-1 rounded-lg">
            Cash: ${user.balance.toFixed(2)}
          </span>
          <span className="bg-slate-800 px-2 py-1 rounded-lg">
            Assets: ${investedCapital.toFixed(2)}
          </span>
        </div>
      </div>

      <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
        <div className="flex items-center gap-3 mb-4">
          <div className="p-2 bg-blue-50 text-blue-600 rounded-lg">
            <Briefcase className="w-6 h-6" />
          </div>
          <div>
            <div className="text-sm text-slate-500 font-bold">Portfolio Items</div>
            <div className="text-2xl font-bold text-slate-900">{portfolioItems.length}</div>
          </div>
        </div>
        <div className="text-sm text-slate-500">Active positions in your portfolio.</div>
      </div>

      <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
        <div className="flex items-center gap-3 mb-4">
          <div className="p-2 bg-green-50 text-green-600 rounded-lg">
            <TrendingUp className="w-6 h-6" />
          </div>
          <div>
            <div className="text-sm text-slate-500 font-bold">Total Gain/Loss</div>
            <div className="text-2xl font-bold text-slate-900">--</div>
          </div>
        </div>
        <div className="text-sm text-slate-500">Real-time P&L not available in this view.</div>
      </div>
    </div>
  );
};
