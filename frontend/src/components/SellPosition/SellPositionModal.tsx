import {
  Modal,
  ModalCard,
  ModalHeader,
  ModalTitle,
  ModalCloseButton,
  ModalFooter,
} from '@/components/Modal';
import { Button } from '@/components/shared/Button';
import { useSellPositionModal } from '@/hooks/useSellPositionModal';
import { SellPositionBody } from './SellPositionBody';

interface SellPositionModalProps {
  isOpen: boolean;
  onClose: () => void;
  symbol: string;
  quantity: number;
  onSuccess?: () => void;
}

export const SellPositionModal = ({
  isOpen,
  onClose,
  symbol,
  quantity,
  onSuccess,
}: SellPositionModalProps) => {
  const {
    displaySymbol,
    price,
    totalValue,
    isPriceLoading,
    isPriceError,
    isTradeLoading,
    handleSellAll,
  } = useSellPositionModal({ isOpen, symbol, quantity, onClose, onSuccess });

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalCard size="sm">
        <ModalHeader>
          <ModalTitle>Sell {displaySymbol}?</ModalTitle>
          <ModalCloseButton />
        </ModalHeader>

        <SellPositionBody
          displaySymbol={displaySymbol}
          quantity={quantity}
          price={price}
          totalValue={totalValue}
          isPriceLoading={isPriceLoading}
          isPriceError={isPriceError}
        />

        <ModalFooter>
          <Button onClick={onClose} variant="outline" className="flex-1">
            Cancel
          </Button>
          <Button
            onClick={handleSellAll}
            disabled={isTradeLoading || isPriceLoading || isPriceError}
            variant="destructive"
            className="flex-1"
          >
            {isTradeLoading ? 'Selling...' : 'Confirm Sell All'}
          </Button>
        </ModalFooter>
      </ModalCard>
    </Modal>
  );
};
