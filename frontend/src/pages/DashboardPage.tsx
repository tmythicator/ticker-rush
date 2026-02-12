import { MarketChart, TradePanel, DashboardStats } from '../components/Dashboard';
import { PortfolioTable } from '../components/PortfolioTable';
import { useQuotesSSE } from '../hooks/useQuotesSSE';
import { useAuth } from '../hooks/useAuth';
import { useTradeSymbol } from '../hooks/useTradeSymbol';

export const DashboardPage = () => {
  const { symbol, setSymbol } = useTradeSymbol();
  const { user } = useAuth();
  const { quote, error: isQuoteError } = useQuotesSSE(symbol);

  return (
    <div className="max-w-[1800px] w-full mx-auto p-4 lg:p-6 grid grid-cols-1 lg:grid-cols-12 gap-6">
      <div className="lg:col-span-9 flex flex-col gap-6">
        <div className="bg-card rounded-lg shadow-sm border border-border p-1 overflow-hidden h-[500px] relative">
          <MarketChart
            key={symbol}
            symbol={symbol}
            onSymbolChange={setSymbol}
            quote={quote || undefined}
            isLoading={!quote}
            isError={!!isQuoteError}
          />
        </div>
        <DashboardStats user={user} />
        <PortfolioTable portfolio={user?.portfolio ?? {}} />
      </div>
      <div className="lg:col-span-3 flex flex-col gap-4 h-full">
        {user && <TradePanel symbol={symbol} currentPrice={quote?.price} />}
      </div>
    </div>
  );
};
