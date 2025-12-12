import { useState } from "react";

export const usePriceColor = (price: number | undefined) => {
    const [color, setColor] = useState('text-slate-900');
    const [prevPrice, setPrevPrice] = useState<number | undefined>(undefined);

    if (price !== undefined && price !== prevPrice) {
        if (prevPrice !== undefined) {
            if (price > prevPrice) {
                setColor('text-emerald-600');
            } else if (price < prevPrice) {
                setColor('text-red-600');
            }
        }
        setPrevPrice(price);
    }

    return color;
};
