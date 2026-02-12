import { useState } from 'react';

export const usePriceColor = (price: number | undefined) => {
  const [color, setColor] = useState('text-foreground');
  const [prevPrice, setPrevPrice] = useState<number | undefined>(undefined);

  if (price !== undefined && price !== prevPrice) {
    if (prevPrice !== undefined) {
      if (price > prevPrice) {
        setColor('text-green-500');
      } else if (price < prevPrice) {
        setColor('text-red-500');
      }
    }
    setPrevPrice(price);
  }

  return color;
};
