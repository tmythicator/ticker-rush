import { useEffect, useRef } from 'react';
import { createChart, ColorType, type ISeriesApi, type Time, AreaSeries } from 'lightweight-charts';
import { type Quote } from '../lib/api';
import { TradeSymbol } from '../types';

interface UseChartProps {
    chartContainerRef: React.RefObject<HTMLDivElement | null>;
    quote?: Quote;
    symbol: TradeSymbol;
}

export const useChart = ({ chartContainerRef, quote, symbol }: UseChartProps) => {
    const chartRef = useRef<ReturnType<typeof createChart> | null>(null);
    const seriesRef = useRef<ISeriesApi<"Area"> | null>(null);

    // Reset data when symbol changes
    useEffect(() => {
        if (seriesRef.current) {
            seriesRef.current.setData([]);
        }
    }, [symbol]);

    // Initialize Chart
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
        //series.setData
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

    // Update Data
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

    return {
        chartRef,
        seriesRef
    };
};
