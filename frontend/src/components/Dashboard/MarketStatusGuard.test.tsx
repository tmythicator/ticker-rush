import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { MarketStatusGuard } from './MarketStatusGuard';

const mockRenderedContent = <div data-testid="market-status-guard-child">Test Rendered Child</div>;

describe('MarketStatusGuard', () => {
  it('renders null when user is not loaded', () => {
    const { container } = render(
      <MarketStatusGuard
        isUserLoaded={false}
        isParticipating={false}
        isLoadingQuotes={false}
        isMarketClosed={false}
      >
        {mockRenderedContent}
      </MarketStatusGuard>,
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders participation guard when user is not participating', () => {
    render(
      <MarketStatusGuard isParticipating={false} isLoadingQuotes={false} isMarketClosed={false}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );

    expect(screen.getByTestId('participation-required-guard')).toBeInTheDocument();
    expect(screen.queryByTestId('market-closed-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('loading-market-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-status-guard-child')).not.toBeInTheDocument();
  });

  it('renders market closed guard when trading is offline', () => {
    render(
      <MarketStatusGuard isParticipating={true} isLoadingQuotes={false} isMarketClosed={true}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );

    expect(screen.getByTestId('market-closed-guard')).toBeInTheDocument();
    expect(screen.queryByTestId('participation-required-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('loading-market-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-status-guard-child')).not.toBeInTheDocument();
  });

  it('renders loading state when quotes are fetching', () => {
    render(
      <MarketStatusGuard isParticipating={true} isLoadingQuotes={true} isMarketClosed={false}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );

    expect(screen.getByTestId('loading-market-guard')).toBeInTheDocument();
    expect(screen.queryByTestId('participation-required-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-closed-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-status-guard-child')).not.toBeInTheDocument();
  });

  it('renders children when all conditions are met', () => {
    render(
      <MarketStatusGuard isParticipating={true} isLoadingQuotes={false} isMarketClosed={false}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );

    expect(screen.getByTestId('market-status-guard-child')).toBeInTheDocument();
    expect(screen.queryByTestId('participation-required-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-closed-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('loading-market-guard')).not.toBeInTheDocument();
  });
});
