export interface TradeFooterProps {
  buyingPower: number;
  estCost: number;
}

export const TradeFooter = ({ buyingPower, estCost }: TradeFooterProps) => {
  return (
    <div className="mt-6 pt-4 border-t border-border">
      <div className="flex justify-between items-center mb-2">
        <span className="text-xs font-bold text-muted-foreground uppercase tracking-wider">
          Buying Power
        </span>
        <span className="font-mono font-bold text-foreground">${buyingPower.toFixed(2)}</span>
      </div>
      <div className="flex justify-between items-center">
        <span className="text-xs font-bold text-muted-foreground uppercase tracking-wider">
          Est. Cost
        </span>
        <span className="font-mono font-bold text-primary">${estCost.toFixed(2)}</span>
      </div>
    </div>
  );
};
