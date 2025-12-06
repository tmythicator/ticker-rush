import { useState } from 'react';

export const usePriceColor = (price: number | undefined) => {
    const [prevPrice, setPrevPrice] = useState<number | null>(null);
    const [priceColor, setPriceColor] = useState('text-slate-900');


    if (price !== undefined && price !== prevPrice) {
        if (prevPrice !== null) {
            if (price > prevPrice) {
                setPriceColor('text-emerald-600');
            } else if (price < prevPrice) {
                setPriceColor('text-red-600');
            }
        }
        setPrevPrice(price);
    }

    return priceColor;
};
