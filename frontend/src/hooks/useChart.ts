import { formatTime, getChartColors } from '@/lib/chartUtils';
import { type Quote, type TradeSymbol } from '@/types';
import {
  AreaSeries,
  ColorType,
  createChart,
  TickMarkType,
  type ChartOptions,
  type DeepPartial,
  type ISeriesApi,
  type UTCTimestamp,
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

    chartRef.current.applyOptions(getThemeOverrides(colors, options));

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
    const container = chartContainerRef.current;
    if (!container) return;

    const colors = getChartColors();
    const chart = createChart(container, {
      localization: { timeFormatter: formatTime },
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
      width: container.clientWidth,
      height: 500,
      autoSize: false,
      timeScale: {
        timeVisible: true,
        secondsVisible: true,
        borderColor: colors.borderColor,
        tickMarkFormatter: formatTickMark,
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
    if (!quote?.timestamp || !seriesRef.current || !symbol) return;

    try {
      const parsedTime = Math.floor(quote.timestamp.getTime() / 1000) as UTCTimestamp;
      seriesRef.current.update({
        time: parsedTime,
        value: quote.price,
      });
    } catch (e) {
      console.warn('Chart update skipped due to timestamp mismatch:', e);
    }
  }, [quote, symbol]);

  return {
    chartRef,
    seriesRef,
  };
};

const formatTickMark = (time: number, tickType: TickMarkType, locale: string): string => {
  const date = new Date(time * 1000);
  const timeOpts: Intl.DateTimeFormatOptions = {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  };

  const formatters: Record<number, () => string> = {
    [TickMarkType.Year]: () => date.getFullYear().toString(),
    [TickMarkType.Month]: () => date.toLocaleString(locale, { month: 'short' }),
    [TickMarkType.DayOfMonth]: () => date.getDate().toString(),
    [TickMarkType.Time]: () => date.toLocaleTimeString(locale, timeOpts),
    [TickMarkType.TimeWithSeconds]: () =>
      date.toLocaleTimeString(locale, { ...timeOpts, second: '2-digit' }),
  };

  return formatters[tickType]?.() ?? date.toLocaleTimeString(locale, timeOpts);
};

const getThemeOverrides = (
  colors: ReturnType<typeof getChartColors>,
  options?: DeepPartial<ChartOptions>,
  // eslint-disable-next-line complexity
) => ({
  layout: {
    background: options?.layout?.background ?? { type: ColorType.Solid, color: colors.bgColor },
    textColor: options?.layout?.textColor ?? colors.textColor,
  },
  grid: {
    horzLines: { color: options?.grid?.horzLines?.color ?? colors.borderColor },
    vertLines: { color: options?.grid?.vertLines?.color ?? colors.borderColor },
  },
  timeScale: {
    borderColor: options?.timeScale?.borderColor ?? colors.borderColor,
  },
  rightPriceScale: {
    borderColor: options?.rightPriceScale?.borderColor ?? colors.borderColor,
  },
});
