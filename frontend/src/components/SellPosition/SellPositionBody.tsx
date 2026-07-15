import { ModalBody } from '@/components/Modal';
import { ErrorMessage } from '@/components/shared/ErrorMessage';
import styles from './SellPosition.module.css';

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
    <ModalBody className={styles.modalBody}>
      <p className={styles.text}>
        Are you sure you want to sell your entire position of{' '}
        <strong className={styles.strongText}>{displaySymbol}</strong>?
      </p>

      <div className={styles.detailsBox}>
        <div className={styles.row}>
          <span className={styles.rowLabel}>Quantity</span>
          <span className={styles.rowValue}>{quantity}</span>
        </div>
        <div className={styles.row}>
          <span className={styles.rowLabel}>Current Price</span>
          {isPriceLoading ? (
            <span className={styles.pulseLoading}>Loading...</span>
          ) : isPriceError ? (
            <span className={styles.unavailable}>Unavailable</span>
          ) : (
            <span className={styles.rowValue}>${price.toFixed(2)}</span>
          )}
        </div>
        <div className={styles.totalRow}>
          <span className={styles.totalLabel}>Total Value</span>
          {isPriceLoading ? (
            <span className={styles.pulseLoading}>Loading...</span>
          ) : isPriceError ? (
            <span className={styles.unavailable}>Unavailable</span>
          ) : (
            <span className={styles.rowValue}>${totalValue.toFixed(2)}</span>
          )}
        </div>
      </div>

      {isPriceError && (
        <ErrorMessage variant="xs">
          Failed to fetch current price. You cannot sell at this time.
        </ErrorMessage>
      )}
    </ModalBody>
  );
};
