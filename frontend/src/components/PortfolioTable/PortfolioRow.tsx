import { IconArrowRight, IconTrash } from '@/components/icons/CustomIcons';
import { SourceBadge } from '@/components/shared/SourceBadge';
import { Button } from '@/components/ui/button';
import { useQuoteQuery } from '@/hooks/useQuoteQuery';
import { formatCurrencyWithSign, parseTicker } from '@/lib/utils';
import type { PortfolioItem, Quote } from '@/types';

interface PortfolioRowProps {
  item: PortfolioItem;
  onSellClick: (item: PortfolioItem) => void;
  onTradeClick: (symbol: string) => void;
  isReadOnly?: boolean;
}

export const PortfolioRow = ({
  item,
  onSellClick,
  onTradeClick,
  isReadOnly = false,
}: PortfolioRowProps) => {
  const { source, symbol } = parseTicker(item.stock_symbol);
  const { data: quote } = useQuoteQuery(item.stock_symbol);

  const isMarketClosed = quote?.is_closed;

  const getPnLColorClass = (q: Quote | undefined) => {
    if (q) {
      return (q.price - item.average_price) * item.quantity >= 0
        ? 'text-green-500'
        : 'text-red-500';
    }

    return 'text-muted-foreground';
  };

  const getPnLElement = (q: Quote | undefined) => {
    if (q) {
      const pnl = (q.price - item.average_price) * item.quantity;
      return <>{formatCurrencyWithSign(pnl)}</>;
    }

    return <span className="animate-pulse">...</span>;
  };

  const getMarketValueElement = (q: Quote | undefined) => {
    if (q) {
      return <>{`$${(q.price * item.quantity).toFixed(2)}`}</>;
    }

    return <span className="animate-pulse">...</span>;
  };

  return (
    <tr key={item.stock_symbol} className="hover:bg-muted/50 transition-colors">
      <td className="px-6 py-4 font-bold text-foreground">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center font-bold text-xs text-muted-foreground">
            {symbol[0].toUpperCase()}
          </div>
          <div className="flex flex-col items-start">
            <span className="text-sm font-bold">{symbol.toUpperCase()}</span>
            <SourceBadge source={source} />
          </div>
        </div>
      </td>
      <td className="px-6 py-4 text-right font-mono text-muted-foreground">{item.quantity}</td>
      <td className="px-6 py-4 text-right font-mono text-muted-foreground">
        ${item.average_price.toFixed(2)}
      </td>
      <td className="px-6 py-4 text-right font-mono text-foreground font-medium">
        {quote ? `$${quote.price.toFixed(2)}` : <span className="animate-pulse">...</span>}
      </td>
      <td className="px-6 py-4 text-right font-mono text-foreground font-bold">
        {getMarketValueElement(quote)}
      </td>
      <td className={`px-6 py-4 text-right font-mono font-bold ${getPnLColorClass(quote)}`}>
        {getPnLElement(quote)}
      </td>
      {!isReadOnly && (
        <td className="px-6 py-4 text-right">
          <div className="flex items-center justify-end gap-2">
            <Button
              variant="destructive"
              size="sm"
              onClick={() => onSellClick(item)}
              title={isMarketClosed ? 'Market Closed' : 'Sell All'}
              className="h-8 px-3 text-xs"
              disabled={isMarketClosed}
            >
              <IconTrash className="w-3 h-3 mr-1" />
              Sell All
            </Button>
            <Button
              variant="default"
              size="sm"
              onClick={() => onTradeClick(item.stock_symbol)}
              className="h-8 px-3 text-xs"
              disabled={isMarketClosed}
              title={isMarketClosed ? 'Market Closed' : 'Trade'}
            >
              Trade
              <IconArrowRight className="w-3 h-3 ml-1" />
            </Button>
          </div>
        </td>
      )}
    </tr>
  );
};
