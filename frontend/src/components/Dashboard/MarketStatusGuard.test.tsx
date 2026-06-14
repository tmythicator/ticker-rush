import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { MarketStatusGuard } from './MarketStatusGuard';
import {
  mockUserParticipating,
  mockUserNotParticipating,
  mockActiveQuote,
  mockClosedQuote,
} from '@/test/mocks';

const mockRenderedContent = <div data-testid="market-status-guard-child">Test Rendered Child</div>;

describe('MarketStatusGuard', () => {
  it('renders nothing when user is null', () => {
    const { container } = render(
      <MarketStatusGuard user={null} quote={mockActiveQuote}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders participation required notice when user is not participating', () => {
    render(
      <MarketStatusGuard user={mockUserNotParticipating} quote={mockActiveQuote}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );

    expect(screen.getByTestId('participation-required-guard')).toBeInTheDocument();
    expect(screen.queryByTestId('market-closed-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('loading-market-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-status-guard-child')).not.toBeInTheDocument();
  });

  it('renders market closed notice when quote is closed', () => {
    render(
      <MarketStatusGuard user={mockUserParticipating} quote={mockClosedQuote}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );

    expect(screen.getByTestId('market-closed-guard')).toBeInTheDocument();
    expect(screen.queryByTestId('participation-required-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('loading-market-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-status-guard-child')).not.toBeInTheDocument();
  });

  it('renders loading indicator when quote is null (e.g. fetcher downtime)', () => {
    render(
      <MarketStatusGuard user={mockUserParticipating} quote={null}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );

    expect(screen.getByTestId('loading-market-guard')).toBeInTheDocument();
    expect(screen.queryByTestId('participation-required-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-closed-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-status-guard-child')).not.toBeInTheDocument();
  });

  it('renders children components when user is participating and market is open', () => {
    render(
      <MarketStatusGuard user={mockUserParticipating} quote={mockActiveQuote}>
        {mockRenderedContent}
      </MarketStatusGuard>,
    );

    expect(screen.getByTestId('market-status-guard-child')).toBeInTheDocument();
    expect(screen.queryByTestId('participation-required-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('market-closed-guard')).not.toBeInTheDocument();
    expect(screen.queryByTestId('loading-market-guard')).not.toBeInTheDocument();
  });
});
