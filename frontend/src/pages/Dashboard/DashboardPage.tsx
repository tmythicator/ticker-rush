import {
  DashboardStats,
  JoinLadderButton,
  MarketChart,
  MarketStatusGuard,
  TradePanel,
} from '@/components/Dashboard';
import { PortfolioHoldings } from '@/components/PortfolioTable';
import { Card } from '@/components/shared/Card';
import { useAuth } from '@/hooks/useAuth';
import { useQuotesSSE } from '@/hooks/useQuotesSSE';
import { useTradeSymbol } from '@/hooks/useTradeSymbol';
import styles from './DashboardPage.module.css';

export const DashboardPage = () => {
  const { symbol, setSymbol } = useTradeSymbol();
  const { user } = useAuth();
  const { quote, error: isQuoteError } = useQuotesSSE(symbol);

  return (
    <div className={styles.dashboardLayout}>
      <h1 className="sr-only">Trading Dashboard</h1>
      <div className={styles.mainSection}>
        <DashboardStats user={user} />

        {!user?.is_participating && (
          <aside className={styles.joinBanner} aria-label="Join Leaderboard">
            <JoinLadderButton />
          </aside>
        )}

        <section className={styles.chartSection} aria-label="Market Chart">
          <Card className={styles.chartCardWrapper}>
            <MarketChart
              key={symbol}
              symbol={symbol}
              onSymbolChange={setSymbol}
              quote={quote}
              isLoading={!quote}
              isError={!!isQuoteError}
            />
          </Card>
        </section>

        <section
          id="trade-panel-mobile"
          className={styles.mobileTradePanel}
          aria-label="Trade Panel"
        >
          <MarketStatusGuard
            isParticipating={user?.is_participating}
            isMarketClosed={quote?.is_closed}
            isLoadingQuotes={!quote}
          >
            <TradePanel quote={quote} />
          </MarketStatusGuard>
        </section>

        <section className={styles.portfolioSection} aria-label="Your Portfolio Holdings">
          <PortfolioHoldings portfolio={user?.portfolio ?? {}} />
        </section>
      </div>

      <aside id="trade-panel-desktop" className={styles.desktopTradePanel} aria-label="Trade Panel">
        <MarketStatusGuard
          isParticipating={user?.is_participating}
          isMarketClosed={quote?.is_closed}
          isLoadingQuotes={!quote}
        >
          <TradePanel quote={quote} />
        </MarketStatusGuard>
      </aside>
    </div>
  );
};
