import { useState } from 'react';

export const usePriceColor = (price: number | undefined) => {
  const [color, setColor] = useState<'neutral' | 'up' | 'down'>('neutral');
  const [prevPrice, setPrevPrice] = useState<number | undefined>(undefined);

  if (price !== undefined && price !== prevPrice) {
    if (prevPrice !== undefined) {
      if (price > prevPrice) {
        setColor('up');
      } else if (price < prevPrice) {
        setColor('down');
      }
    }
    setPrevPrice(price);
  }

  return color;
};
