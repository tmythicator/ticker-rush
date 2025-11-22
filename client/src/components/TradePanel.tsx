import { TrendingUp } from 'lucide-react';

const TradeButton = ({ type, onClick }: { type: 'BUY' | 'SELL', onClick?: () => void }) => {
    const isBuy = type === 'BUY';

    const baseStyles = "w-full font-bold py-3 rounded-lg transition-all shadow-sm active:scale-95 transform duration-100 flex items-center justify-center gap-2";
    const activeStyles = isBuy
        ? "bg-green-600 hover:bg-green-700 text-white shadow-green-100"
        : "bg-red-600 hover:bg-red-700 text-white shadow-red-100";

    return (
        <button
            onClick={onClick}
            className={`${baseStyles} ${activeStyles}`}
        >
            {type}
        </button>
    );
};

export const TradePanel = () => {
    return (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 flex flex-col h-full">
            <div className="flex items-center justify-between mb-6">
                <h2 className="font-bold text-slate-800 flex items-center gap-2">
                    <TrendingUp className="w-4 h-4 text-blue-600" />
                    Place Order
                </h2>
            </div>

            <div className="space-y-5 flex-1">
                <div>
                    <label className="block text-xs font-bold text-slate-400 mb-2 uppercase tracking-wider">Symbol</label>
                    <div className="relative">
                        <input
                            type="text"
                            value="AAPL"
                            disabled
                            className="w-full bg-slate-50 border border-slate-200 rounded-lg px-4 py-3 font-mono text-sm font-bold text-slate-700 opacity-70"
                        />
                        <div className="absolute right-3 top-3 text-xs font-bold text-slate-400">STOCK</div>
                    </div>
                </div>

                <div>
                    <label className="block text-xs font-bold text-slate-400 mb-2 uppercase tracking-wider">Quantity</label>
                    <div className="relative">
                        <input
                            type="number"
                            placeholder="0"
                            className="w-full bg-white border border-slate-200 rounded-lg px-4 py-3 font-mono text-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all shadow-sm placeholder:text-slate-300"
                        />
                    </div>
                </div>

                <div className="pt-2 grid grid-cols-2 gap-3">
                    <TradeButton type="BUY" />
                    <TradeButton type="SELL" />
                </div>
            </div>

            <div className="mt-auto pt-6 border-t border-slate-100 space-y-2">
                <div className="flex justify-between text-sm">
                    <span className="text-slate-500">Buying Power</span>
                    <span className="font-mono font-medium text-slate-700">$4,250.00</span>
                </div>
                <div className="flex justify-between text-sm">
                    <span className="text-slate-500">Est. Cost</span>
                    <span className="font-mono font-bold text-slate-900">$0.00</span>
                </div>
            </div>
        </div>
    );
};