export interface TradeFooterProps {
    buyingPower: number;
    estCost: number;
}

export const TradeFooter = ({ buyingPower, estCost }: TradeFooterProps) => {
    const formatCurrency = (val: number) =>
        new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);
    return (
        <div className="mt-auto pt-6 border-t border-slate-100 space-y-2">
            <div className="flex justify-between text-sm">
                <span className="text-slate-500">Buying Power</span>
                <span className="font-mono font-medium text-slate-700">{formatCurrency(buyingPower)}</span>
            </div>
            <div className="flex justify-between text-sm">
                <span className="text-slate-500">Est. Cost</span>
                <span className="font-mono font-bold text-slate-900">{formatCurrency(estCost)}</span>
            </div>
        </div>);
};