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
    <div className={`flex-1 space-y-5 ${disabled ? 'pointer-events-none opacity-50' : ''}`}>
      {error && <div className="mb-2 text-xs font-bold text-red-600">{error}</div>}

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
