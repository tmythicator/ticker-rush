import { Trophy, Medal, TrendingUp, User } from 'lucide-react';

// TODO: get real data from backend
const leaders = [
    { rank: 1, name: "Alex. B", profit: "+245.0%", volume: "$12.5M" },
    { rank: 2, name: "Tim R.", profit: "+180.2%", volume: "$8.2M" },
    { rank: 3, name: "Garry O.", profit: "+150.5%", volume: "$5.1M" },
    { rank: 4, name: "Timmy", profit: "+98.0%", volume: "$2.1M" },
    { rank: 5, name: "Jonas B.", profit: "+45.2%", volume: "$1.8M" },
];

export const Leaderboard = () => {
    return (
        <div className="max-w-5xl w-full mx-auto p-4 lg:p-6">
            <div className="mb-8 text-center">
                <h2 className="text-3xl font-bold text-slate-900 flex items-center justify-center gap-3">
                    <Trophy className="w-8 h-8 text-yellow-500" />
                    Season 1 Ladder
                </h2>
                <p className="text-slate-500 mt-2">Top traders by monthly margin</p>
            </div>

            <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm">
                        <thead className="bg-slate-50 border-b border-slate-200">
                            <tr>
                                <th className="px-6 py-4 font-bold text-slate-500 uppercase tracking-wider w-20">Rank</th>
                                <th className="px-6 py-4 font-bold text-slate-500 uppercase tracking-wider">Trader</th>
                                <th className="px-6 py-4 font-bold text-slate-500 uppercase tracking-wider text-right">Volume</th>
                                <th className="px-6 py-4 font-bold text-slate-500 uppercase tracking-wider text-right">Total P&L</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-slate-100">
                            {leaders.map((leader) => {
                                let rankIcon = <span className="font-mono font-bold text-slate-500">#{leader.rank}</span>;
                                let rowClass = "hover:bg-slate-50 transition-colors";

                                if (leader.rank === 1) {
                                    rankIcon = <Trophy className="w-5 h-5 text-yellow-500" />;
                                    rowClass = "bg-yellow-50/30 hover:bg-yellow-50/50";
                                } else if (leader.rank === 2) {
                                    rankIcon = <Medal className="w-5 h-5 text-slate-400" />;
                                } else if (leader.rank === 3) {
                                    rankIcon = <Medal className="w-5 h-5 text-amber-700" />;
                                }

                                return (
                                    <tr key={leader.rank} className={rowClass}>
                                        <td className="px-6 py-4 flex justify-center">{rankIcon}</td>
                                        <td className="px-6 py-4">
                                            <div className="flex items-center gap-3">
                                                <div className="w-8 h-8 bg-slate-100 rounded-full flex items-center justify-center border border-slate-200">
                                                    <User className="w-4 h-4 text-slate-400" />
                                                </div>
                                                <div>
                                                    <span className="font-bold text-slate-900 block">{leader.name}</span>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-right font-mono text-slate-600">{leader.volume}</td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="font-bold text-green-600 flex items-center justify-end gap-1">
                                                <TrendingUp className="w-3 h-3" />
                                                {leader.profit}
                                            </div>
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};