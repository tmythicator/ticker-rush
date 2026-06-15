import {
  DashboardStats,
  JoinLadderButton,
  MarketChart,
  MarketStatusGuard,
  TradePanel,
} from '@/components/Dashboard';
import { PortfolioHoldings } from '@/components/PortfolioTable';
import { useAuth } from '@/hooks/useAuth';
import { useQuotesSSE } from '@/hooks/useQuotesSSE';
import { useTradeSymbol } from '@/hooks/useTradeSymbol';

export const DashboardPage = () => {
  const { symbol, setSymbol } = useTradeSymbol();
  const { user } = useAuth();
  const { quote, error: isQuoteError } = useQuotesSSE(symbol);

  return (
    <div className="mx-auto grid w-full max-w-[1800px] grid-cols-1 gap-6 p-4 pb-6 lg:grid-cols-12 lg:p-6">
      <div className="flex flex-col gap-6 lg:col-span-9">
        <DashboardStats user={user} />

        {!user?.is_participating && (
          <div className="duration-500 animate-in fade-in slide-in-from-top-4">
            <JoinLadderButton />
          </div>
        )}

        <div className="relative h-[500px] overflow-hidden rounded-lg border border-border bg-card p-1 shadow-sm">
          <MarketChart
            key={symbol}
            symbol={symbol}
            onSymbolChange={setSymbol}
            quote={quote}
            isLoading={!quote}
            isError={!!isQuoteError}
          />
        </div>

        <div className="lg:hidden" id="trade-panel-mobile">
          <MarketStatusGuard user={user} quote={quote}>
            <TradePanel quote={quote} />
          </MarketStatusGuard>
        </div>

        <PortfolioHoldings portfolio={user?.portfolio ?? {}} />
      </div>

      <div className="hidden h-full flex-col gap-4 lg:col-span-3 lg:flex" id="trade-panel-desktop">
        <MarketStatusGuard user={user} quote={quote}>
          <TradePanel quote={quote} />
        </MarketStatusGuard>
      </div>
    </div>
  );
};
