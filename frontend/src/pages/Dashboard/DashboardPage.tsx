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
      <div className={styles.mainSection}>
        <DashboardStats user={user} />

        {!user?.is_participating && (
          <div className={styles.joinBanner}>
            <JoinLadderButton />
          </div>
        )}

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

        <div className={styles.mobileTradePanel} id="trade-panel-mobile">
          <MarketStatusGuard user={user} quote={quote}>
            <TradePanel quote={quote} />
          </MarketStatusGuard>
        </div>

        <PortfolioHoldings portfolio={user?.portfolio ?? {}} />
      </div>

      <div className={styles.desktopTradePanel} id="trade-panel-desktop">
        <MarketStatusGuard user={user} quote={quote}>
          <TradePanel quote={quote} />
        </MarketStatusGuard>
      </div>
    </div>
  );
};
