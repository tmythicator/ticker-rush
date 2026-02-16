import { getChartColors } from '@/lib/chartUtils';
import { type Quote, type TradeSymbol } from '@/types';
import {
  AreaSeries,
  ColorType,
  createChart,
  type ChartOptions,
  type DeepPartial,
  type ISeriesApi,
  type Time,
} from 'lightweight-charts';
import { useEffect, useLayoutEffect, useRef } from 'react';
import { useChartHistory } from './useChartHistory';
import { useThemeObserver } from './useThemeObserver';

interface UseChartProps {
  chartContainerRef: React.RefObject<HTMLDivElement | null>;
  quote: Quote | null;
  symbol: TradeSymbol | null;
  options?: DeepPartial<ChartOptions>;
}

export const useChart = ({ chartContainerRef, quote, symbol, options }: UseChartProps) => {
  const chartRef = useRef<ReturnType<typeof createChart> | null>(null);
  const seriesRef = useRef<ISeriesApi<'Area'> | null>(null);

  useChartHistory(symbol, seriesRef);

  const updateChartColors = () => {
    if (!chartRef.current) return;
    const colors = getChartColors();

    chartRef.current.applyOptions({
      layout: {
        background: options?.layout?.background || { type: ColorType.Solid, color: colors.bgColor },
        textColor: options?.layout?.textColor || colors.textColor,
      },
      grid: {
        horzLines: { color: options?.grid?.horzLines?.color || colors.borderColor },
        vertLines: { color: options?.grid?.vertLines?.color || colors.borderColor },
      },
      timeScale: {
        borderColor: options?.timeScale?.borderColor || colors.borderColor,
      },
      rightPriceScale: {
        borderColor: options?.rightPriceScale?.borderColor || colors.borderColor,
      },
    });

    if (seriesRef.current) {
      seriesRef.current.applyOptions({
        topColor: colors.areaTopColor,
        bottomColor: colors.areaBottomColor,
        lineColor: colors.areaLineColor,
      });
    }
  };

  useThemeObserver(updateChartColors);

  // Initialize Chart
  useLayoutEffect(() => {
    if (!chartContainerRef.current) return;

    const colors = getChartColors();
    const chart = createChart(chartContainerRef.current, {
      localization: {
        timeFormatter: (timestamp: number) => {
          return new Date(timestamp * 1000).toLocaleTimeString(undefined, {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false,
          });
        },
      },
      ...options,
      layout: {
        background: { type: ColorType.Solid, color: colors.bgColor },
        textColor: colors.textColor,
        attributionLogo: false,
        ...options?.layout,
      },
      grid: {
        vertLines: { visible: false },
        horzLines: { color: colors.borderColor },
        ...options?.grid,
      },
      width: chartContainerRef.current.clientWidth,
      height: 500,
      autoSize: true,
      timeScale: {
        timeVisible: true,
        secondsVisible: true,
        borderColor: colors.borderColor,
        ...options?.timeScale,
      },
      rightPriceScale: {
        borderColor: colors.borderColor,
        ...options?.rightPriceScale,
      },
    });

    chartRef.current = chart;

    const series = chart.addSeries(AreaSeries, {
      topColor: colors.areaTopColor,
      bottomColor: colors.areaBottomColor,
      lineColor: colors.areaLineColor,
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

      seriesRef.current = null;
    };
  }, [chartContainerRef, options]);

  // Update Data
  useEffect(() => {
    if (quote && seriesRef.current && symbol) {
      try {
        seriesRef.current.update({
          time: quote.timestamp as Time,
          value: quote.price,
        });
      } catch (e) {
        console.warn('Chart update skipped due to timestamp mismatch:', e);
      }
    }
  }, [quote, symbol]);

  return {
    chartRef,
    seriesRef,
  };
};
