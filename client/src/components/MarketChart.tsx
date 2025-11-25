import { useEffect, useRef } from 'react';
import { createChart, ColorType, CandlestickSeries, type ISeriesApi, type Time } from 'lightweight-charts';
import { useQuery } from '@tanstack/react-query';
import { fetchQuote } from '../lib/api';

export const MarketChart = () => {
    const chartContainerRef = useRef<HTMLDivElement>(null);
    const seriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);

    const { data: quote } = useQuery({
        queryKey: ['quote', 'AAPL'],
        queryFn: () => fetchQuote('AAPL'),
        refetchInterval: 2000,
    });

    useEffect(() => {
        if (!chartContainerRef.current) return;

        const chart = createChart(chartContainerRef.current, {
            layout: {
                background: { type: ColorType.Solid, color: 'white' },
                textColor: '#333',
                attributionLogo: false,
            },
            grid: {
                vertLines: { color: '#f0f3fa' },
                horzLines: { color: '#f0f3fa' },
            },
            width: chartContainerRef.current.clientWidth,
            height: 500,
            autoSize: true,
            timeScale: { timeVisible: true, secondsVisible: true }
        });

        const series = chart.addSeries(CandlestickSeries, {
            upColor: '#26a69a', downColor: '#ef5350', borderVisible: false,
            wickUpColor: '#26a69a', wickDownColor: '#ef5350',
        });
        seriesRef.current = series;

        chart.timeScale().fitContent();

        return () => {
            chart.remove();
        };
    }, []);

    useEffect(() => {
        if (quote && seriesRef.current) {
            seriesRef.current.update({
                time: quote.timestamp as Time,
                open: quote.price,
                high: quote.price + 0.5,
                low: quote.price - 0.5,
                close: quote.price,
            });
        }
    }, [quote]);


    const price = quote?.price;
    const isUp = price && price > 150;


    return (
        <div className="w-full h-full relative group">
            <div ref={chartContainerRef} className="w-full h-[500px]" />
            <div className="absolute top-3 left-3 bg-white/90 backdrop-blur-sm px-3 py-1.5 rounded-md shadow-sm border border-slate-200 z-10 flex items-center gap-2">
                <span className="font-bold text-slate-800 tracking-tight">AAPL</span>
                <span className="text-xs text-slate-500 font-medium border-l border-slate-300 pl-2">Apple Inc.</span>
                {price && (
                    <span className={`text-sm font-mono font-bold transition-colors ${isUp ? 'text-green-600' : 'text-red-600'}`}>
                        ${price.toFixed(2)}
                    </span>
                )}
            </div>
        </div>
    );
};