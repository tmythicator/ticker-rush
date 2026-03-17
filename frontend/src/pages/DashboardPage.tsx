import {
  DashboardStats,
  JoinLadderButton,
  MarketChart,
  MarketStatusGuard,
  TradePanel,
} from '@/components/Dashboard';
import { PortfolioTable } from '@/components/PortfolioTable/PortfolioTable';
import { useAuth } from '@/hooks/useAuth';
import { useQuotesSSE } from '@/hooks/useQuotesSSE';
import { useTradeSymbol } from '@/hooks/useTradeSymbol';

export const DashboardPage = () => {
  const { symbol, setSymbol } = useTradeSymbol();
  const { user } = useAuth();
  const { quote, error: isQuoteError } = useQuotesSSE(symbol);

  return (
    <div className="max-w-[1800px] w-full mx-auto p-4 lg:p-6 grid grid-cols-1 lg:grid-cols-12 gap-6 pb-6">
      <div className="lg:col-span-9 flex flex-col gap-6">
        <DashboardStats user={user} />

        {!user?.is_participating && (
          <div className="animate-in fade-in slide-in-from-top-4 duration-500">
            <JoinLadderButton />
          </div>
        )}

        <div className="bg-card rounded-lg shadow-sm border border-border p-1 overflow-hidden h-[500px] relative">
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

        <PortfolioTable portfolio={user?.portfolio ?? {}} />
      </div>

      <div className="hidden lg:flex lg:col-span-3 flex-col gap-4 h-full" id="trade-panel-desktop">
        <MarketStatusGuard user={user} quote={quote}>
          <TradePanel quote={quote} />
        </MarketStatusGuard>
      </div>
    </div>
  );
};
