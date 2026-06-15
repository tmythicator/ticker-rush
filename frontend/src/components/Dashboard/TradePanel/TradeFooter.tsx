export interface TradeFooterProps {
  buyingPower: number;
  estCost: number;
}

export const TradeFooter = ({ buyingPower, estCost }: TradeFooterProps) => {
  return (
    <div className="mt-6 border-t border-border pt-4">
      <div className="mb-2 flex items-center justify-between">
        <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
          Buying Power
        </span>
        <span className="font-mono font-bold text-foreground">${buyingPower.toFixed(2)}</span>
      </div>
      <div className="flex items-center justify-between">
        <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
          Est. Cost
        </span>
        <span className="font-mono font-bold text-primary">${estCost.toFixed(2)}</span>
      </div>
    </div>
  );
};
