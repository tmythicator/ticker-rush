import { useEffect, useRef, useState } from 'react';
import { createChart, ColorType, type ISeriesApi, type Time, AreaSeries } from 'lightweight-charts';
import { useQuery } from '@tanstack/react-query';
import { fetchQuote } from '../lib/api';
import { ChevronDown } from 'lucide-react';

const TICKERS = ["AAPL", "BINANCE:BTCUSDT"];

export const MarketChart = () => {
    const chartContainerRef = useRef<HTMLDivElement>(null);
    const seriesRef = useRef<ISeriesApi<"Area"> | null>(null);
    const chartRef = useRef<ReturnType<typeof createChart> | null>(null);
    const prevPriceRef = useRef<number | null>(null);

    const [symbol, setSymbol] = useState("AAPL");

    const { data: quote, isLoading, isError } = useQuery({
        queryKey: ['quote', symbol],
        queryFn: () => fetchQuote(symbol),
        refetchInterval: 3000,
    });

    useEffect(() => {
        if (seriesRef.current) {
            seriesRef.current.setData([]);
            prevPriceRef.current = null;
        }
    }, [symbol]);

    useEffect(() => {
        if (!chartContainerRef.current) return;

        const chart = createChart(chartContainerRef.current, {
            layout: {
                background: { type: ColorType.Solid, color: 'white' },
                attributionLogo: false,
                textColor: '#333',
            },
            grid: {
                vertLines: { visible: false },
                horzLines: { color: '#f0f3fa' }
            },
            width: chartContainerRef.current.clientWidth,
            height: 500,
            autoSize: true,
            timeScale: {
                timeVisible: true,
                secondsVisible: true,
                borderColor: '#f0f3fa',
            },
            rightPriceScale: {
                borderColor: '#f0f3fa',
            }
        });

        chartRef.current = chart;

        const series = chart.addSeries(AreaSeries, {
            topColor: 'rgba(38, 166, 154, 0.56)',
            bottomColor: 'rgba(38, 166, 154, 0.04)',
            lineColor: 'rgba(38, 166, 154, 1)',
            lineWidth: 2,
        });
        seriesRef.current = series;

        const handleResize = () => {
            if (chartContainerRef.current) {
                chart.applyOptions({ width: chartContainerRef.current.clientWidth });
            }
        };
        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('resize', handleResize);
            chart.remove();
        };
    }, []);

    useEffect(() => {
        if (quote && seriesRef.current) {
            try {
                seriesRef.current.update({
                    time: quote.timestamp as Time,
                    value: quote.price,
                });
            } catch (e) {
                console.warn("Chart update skipped due to timestamp mismatch:", e);
            }
        }
    }, [quote, symbol]);

    const price = quote?.price;
    let priceColor = 'text-slate-900';

    if (price !== undefined && prevPriceRef.current !== null) {
        if (price > prevPriceRef.current) priceColor = 'text-emerald-600';
        else if (price < prevPriceRef.current) priceColor = 'text-red-600';
    }


    useEffect(() => {
        if (price !== undefined) {
            prevPriceRef.current = price;
        }
    }, [price]);

    return (
        <div className="w-full h-full relative group">
            <div ref={chartContainerRef} className="w-full h-[500px]" />
            <div className="absolute top-4 left-4 z-20 flex flex-col gap-2">
                <div className="bg-white/90 backdrop-blur-md p-1.5 rounded-xl shadow-lg border border-slate-200/60 flex items-center gap-1 transition-all hover:shadow-xl hover:scale-[1.02]">
                    <div className="relative group/select">
                        <select
                            value={symbol}
                            onChange={(e) => setSymbol(e.target.value)}
                            className="appearance-none bg-transparent pl-3 pr-8 py-1.5 font-bold text-slate-800 text-lg tracking-tight focus:outline-none cursor-pointer hover:text-blue-600 transition-colors"
                        >
                            {TICKERS.map((t) => (
                                <option key={t} value={t}>{t}</option>
                            ))}
                        </select>
                        <ChevronDown className="w-4 h-4 text-slate-400 absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none group-hover/select:text-blue-500 transition-colors" />
                    </div>

                    <div className="h-6 w-px bg-slate-200 mx-1"></div>

                    <div className="px-2 flex flex-col items-end min-w-[80px]">
                        {isLoading ? (
                            <div className="h-5 w-16 bg-slate-200 animate-pulse rounded"></div>
                        ) : isError ? (
                            <span className="text-xs font-bold text-red-500">OFFLINE</span>
                        ) : (
                            <>
                                <span className={`text-lg font-mono font-bold leading-none ${priceColor} transition-colors duration-300`}>
                                    {price ? `$${price.toFixed(2)}` : '—'}
                                </span>
                                <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Live</span>
                            </>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};