import { TradeAction, type TickerSource } from '@/types';
import { TradeButtons } from './TradeButtons';
import { SymbolField } from './SymbolField';
import { QuantityField } from './QuantityField';

export interface TradeOrderAsset {
  symbol: string;
  source?: TickerSource;
  price?: number;
  positionQuantity?: number;
  buyingPower?: number;
}

export interface TradeOrderFormState {
  quantity: string;
  setQuantity: (quantity: string) => void;
  error: string | null;
  disabled?: boolean;
}

interface TradeOrderInputProps {
  asset: TradeOrderAsset;
  form: TradeOrderFormState;
  onTrade: (action: TradeAction) => void;
}

export const TradeOrderInput = ({ asset, form, onTrade }: TradeOrderInputProps) => {
  const { symbol, source, price, positionQuantity, buyingPower } = asset;
  const { quantity, setQuantity, error, disabled } = form;

  return (
    <div className={`space-y-5 flex-1 ${disabled ? 'opacity-50 pointer-events-none' : ''}`}>
      {error && <div className="text-xs text-red-600 font-bold mb-2">{error}</div>}

      <SymbolField symbol={symbol} source={source} />

      <QuantityField
        quantity={quantity}
        setQuantity={setQuantity}
        buyingPower={buyingPower}
        price={price}
        disabled={disabled}
        positionQuantity={positionQuantity}
      />

      <TradeButtons handleTrade={onTrade} disabled={disabled} />
    </div>
  );
};
