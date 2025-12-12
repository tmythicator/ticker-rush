import { TrendingUp, RefreshCcw } from "lucide-react";

export const TradePanelHeader = ({ isLoading }: { isLoading: boolean }) => {
    return (
        <div className="flex items-center justify-between mb-6">
            <h2 className="font-bold text-slate-800 flex items-center gap-2">
                <TrendingUp className="w-4 h-4 text-blue-600" />
                Place Order
            </h2>
            {isLoading && <RefreshCcw className="w-4 h-4 animate-spin text-slate-400" />}
        </div>
    );
}