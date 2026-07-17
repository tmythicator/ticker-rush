import type { TradeAction } from '@/types';
import { type TickerSource } from '@/types';
import { TradeButtons } from './TradeButtons';
import { SymbolField } from './SymbolField';
import { QuantityField } from './QuantityField';
import { ErrorMessage } from '@/components/shared/ErrorMessage';
import styles from './TradePanel.module.css';

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
    <div className={styles.formContainer} data-disabled={disabled}>
      {error && <ErrorMessage variant="xs" message={error} />}

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
