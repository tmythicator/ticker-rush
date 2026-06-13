import { ModalBody } from '@/components/Modal';

interface SellPositionBodyProps {
  displaySymbol: string;
  quantity: number;
  price: number;
  totalValue: number;
  isPriceLoading: boolean;
  isPriceError: boolean;
}

export const SellPositionBody = ({
  displaySymbol,
  quantity,
  price,
  totalValue,
  isPriceLoading,
  isPriceError,
}: SellPositionBodyProps) => {
  return (
    <ModalBody className="space-y-4">
      <p className="text-muted-foreground text-sm">
        Are you sure you want to sell your entire position of{' '}
        <strong className="text-foreground">{displaySymbol}</strong>?
      </p>

      <div className="bg-muted/50 rounded-lg p-4 space-y-2 border border-border">
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">Quantity</span>
          <span className="font-mono font-bold text-foreground">{quantity}</span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">Current Price</span>
          {isPriceLoading ? (
            <span className="text-muted-foreground animate-pulse">Loading...</span>
          ) : isPriceError ? (
            <span className="text-destructive font-bold">Unavailable</span>
          ) : (
            <span className="font-mono font-bold text-foreground">${price.toFixed(2)}</span>
          )}
        </div>
        <div className="border-t border-border pt-2 flex justify-between text-sm">
          <span className="font-bold text-foreground">Total Value</span>
          {isPriceLoading ? (
            <span className="text-muted-foreground animate-pulse">Loading...</span>
          ) : isPriceError ? (
            <span className="text-destructive font-bold">Unavailable</span>
          ) : (
            <span className="font-mono font-bold text-foreground">${totalValue.toFixed(2)}</span>
          )}
        </div>
      </div>

      {isPriceError && (
        <div className="text-xs text-destructive bg-destructive/10 p-2 rounded-lg border border-destructive/20">
          Failed to fetch current price. You cannot sell at this time.
        </div>
      )}
    </ModalBody>
  );
};
