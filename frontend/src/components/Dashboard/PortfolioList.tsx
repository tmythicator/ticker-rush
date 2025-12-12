import { useSearchParams } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

export const PortfolioList = () => {
    const { user } = useAuth();
    const positions = user?.portfolio ? Object.values(user.portfolio) : [];
    const [, setSearchParams] = useSearchParams();

    const handleRowClick = (symbol: string) => {
        setSearchParams({ symbol });
    };

    if (positions.length === 0) {
        return (
            <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
                <h3 className="text-lg font-bold text-slate-900 mb-4">Your Positions</h3>
                <div className="text-slate-500 text-center py-8">
                    No open positions. Start trading!
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
            <h3 className="text-lg font-bold text-slate-900 mb-4">Your Positions</h3>
            <div className="overflow-x-auto">
                <table className="w-full text-left font-medium">
                    <thead>
                        <tr className="border-b border-slate-100 text-slate-500 text-sm">
                            <th className="pb-3 pl-2">Symbol</th>
                            <th className="pb-3 text-right">Shares</th>
                            <th className="pb-3 text-right">Avg Price</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-50">
                        {positions.map((item) => (
                            <tr
                                key={item.stock_symbol}
                                onClick={() => handleRowClick(item.stock_symbol)}
                                className="group hover:bg-slate-50 transition-colors cursor-pointer"
                            >
                                <td className="py-3 pl-2 text-slate-900 font-bold">{item.stock_symbol}</td>
                                <td className="py-3 text-right text-slate-700">{item.quantity}</td>
                                <td className="py-3 text-right text-slate-500">${item.average_price.toFixed(2)}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};
