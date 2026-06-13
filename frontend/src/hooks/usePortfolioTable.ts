import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { PortfolioItem } from '@/types';

export const usePortfolioTable = (portfolio: Record<string, PortfolioItem>) => {
  const navigate = useNavigate();
  const [sellModal, setSellModal] = useState<{ isOpen: boolean; item?: PortfolioItem }>({
    isOpen: false,
  });

  const handleTrade = (symbol: string) => {
    navigate(`/trade?symbol=${symbol}`);
  };

  const handleSellClick = (item: PortfolioItem) => {
    setSellModal({ isOpen: true, item });
  };

  const handleCloseSellModal = () => {
    setSellModal({ isOpen: false });
  };

  const items = Object.values(portfolio);
  const isEmpty = items.length === 0;

  return {
    items,
    isEmpty,
    sellModal,
    handleTrade,
    handleSellClick,
    handleCloseSellModal,
  };
};
