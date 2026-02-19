import { DashboardStats, MarketChart, TradePanel } from '@/components/Dashboard';
import { IconMoon } from '@/components/icons/CustomIcons';
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

        {/* Mobile Trade Panel (Between Chart and Portfolio) */}
        <div className="lg:hidden" id="trade-panel-mobile">
          {user &&
            (quote?.is_closed ? (
              <div className="bg-card rounded-lg shadow-sm border border-border p-6 flex flex-col items-center justify-center text-center">
                <IconMoon className="w-8 h-8 mb-4 text-primary" />
                <h3 className="text-xl font-bold text-foreground">Market Closed</h3>
                <p className="text-muted-foreground mt-2">
                  Trading is currently unavailable.
                  <br />
                  Please come back during market hours.
                </p>
              </div>
            ) : (
              <TradePanel quote={quote} />
            ))}
        </div>

        <PortfolioTable portfolio={user?.portfolio ?? {}} />
      </div>
      <div className="hidden lg:flex lg:col-span-3 flex-col gap-4 h-full" id="trade-panel-desktop">
        {user &&
          (quote?.is_closed ? (
            <div className="bg-card rounded-lg shadow-sm border border-border p-6 flex flex-col items-center justify-center text-center h-full">
              <IconMoon className="w-8 h-8 mb-4 text-primary" />

              <h3 className="text-xl font-bold text-foreground">Market Closed</h3>
              <p className="text-muted-foreground mt-2">
                Trading is currently unavailable.
                <br />
                Please come back during market hours.
              </p>
            </div>
          ) : (
            <TradePanel quote={quote} />
          ))}
      </div>
    </div>
  );
};
