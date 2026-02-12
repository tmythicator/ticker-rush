import { type Quote } from '@/lib/api';
import { getChartColors } from '@/lib/chart-utils';
import { AreaSeries, ColorType, createChart, type ISeriesApi, type Time } from 'lightweight-charts';
import { useEffect, useRef } from 'react';
import { TradeSymbol } from '../types';
import { useThemeObserver } from './useThemeObserver';

interface UseChartProps {
  chartContainerRef: React.RefObject<HTMLDivElement | null>;
  quote?: Quote;
  symbol: TradeSymbol;
}

export const useChart = ({ chartContainerRef, quote, symbol }: UseChartProps) => {
  const chartRef = useRef<ReturnType<typeof createChart> | null>(null);
  const seriesRef = useRef<ISeriesApi<'Area'> | null>(null);

  // Reset data when symbol changes
  useEffect(() => {
    if (seriesRef.current) {
      seriesRef.current.setData([]);
    }
  }, [symbol]);

  const updateChartColors = () => {
    if (!chartRef.current) return;
    const colors = getChartColors();

    chartRef.current.applyOptions({
      layout: {
        background: { type: ColorType.Solid, color: colors.bgColor },
        textColor: colors.textColor,
      },
      grid: {
        horzLines: { color: colors.borderColor },
        vertLines: { color: colors.borderColor },
      },
      timeScale: {
        borderColor: colors.borderColor,
      },
      rightPriceScale: {
        borderColor: colors.borderColor,
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
  useEffect(() => {
    if (!chartContainerRef.current) return;

    const colors = getChartColors();
    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: colors.bgColor },
        textColor: colors.textColor,
        attributionLogo: false,
      },
      grid: {
        vertLines: { visible: false },
        horzLines: { color: colors.borderColor },
      },
      width: chartContainerRef.current.clientWidth,
      height: 500,
      autoSize: true,
      timeScale: {
        timeVisible: true,
        secondsVisible: true,
        borderColor: colors.borderColor,
      },
      rightPriceScale: {
        borderColor: colors.borderColor,
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
    };
  }, [chartContainerRef]);

  // Update Data
  useEffect(() => {
    if (quote && seriesRef.current) {
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
