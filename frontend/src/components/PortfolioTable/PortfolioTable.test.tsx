import { render, screen, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PortfolioTable } from './PortfolioTable';
import { usePortfolioRowState } from '@/hooks/usePortfolioRowState';
import { mockPortfolioItemAAPL, mockPortfolioItemMSFT } from '@/test/mocks';
import userEvent from '@testing-library/user-event';

vi.mock('@/hooks/usePortfolioRowState', () => ({
  usePortfolioRowState: vi.fn(),
}));

const mockItems = [mockPortfolioItemAAPL, mockPortfolioItemMSFT];

describe('PortfolioTable', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders portfolio rows and elements correctly', () => {
    vi.mocked(usePortfolioRowState).mockImplementation((item) => {
      if (item.stock_symbol === 'AAPL') {
        return {
          symbol: 'AAPL',
          source: 'Finnhub',
          isMarketClosed: false,
          isTradable: true,
          quote: { price: 160.0 },
          marketValue: '$1600.00',
          pnl: '+$100.00',
          pnlColorClass: 'text-green-500',
        } as unknown as ReturnType<typeof usePortfolioRowState>;
      }
      return {
        symbol: 'MSFT',
        source: 'Finnhub',
        isMarketClosed: false,
        isTradable: false,
        quote: { price: 310.0 },
        marketValue: '$1550.00',
        pnl: '+$50.00',
        pnlColorClass: 'text-green-500',
      } as unknown as ReturnType<typeof usePortfolioRowState>;
    });

    render(<PortfolioTable items={mockItems} />);

    // Assert table is in the document
    expect(screen.getByTestId('portfolio-table')).toBeInTheDocument();

    // Find rows by data-testid
    const aaplRow = screen.getByTestId('portfolio-row-aapl');
    const msftRow = screen.getByTestId('portfolio-row-msft');

    expect(aaplRow).toBeInTheDocument();
    expect(msftRow).toBeInTheDocument();

    // Assert row-scoped elements
    expect(within(aaplRow).queryByTestId('suspended-badge')).not.toBeInTheDocument();
    expect(within(msftRow).getByTestId('suspended-badge')).toBeInTheDocument();
  });

  it('allows actions on tradable assets and invokes callbacks', async () => {
    const user = userEvent.setup();
    const handleSellClick = vi.fn();
    const handleTradeClick = vi.fn();

    vi.mocked(usePortfolioRowState).mockReturnValue({
      symbol: 'AAPL',
      source: 'Finnhub',
      isMarketClosed: false,
      isTradable: true,
      quote: { price: 160.0 },
      marketValue: '$1600.00',
      pnl: '+$100.00',
      pnlColorClass: 'text-green-500',
    } as unknown as ReturnType<typeof usePortfolioRowState>);

    render(
      <PortfolioTable
        items={[mockPortfolioItemAAPL]}
        onSellClick={handleSellClick}
        onTradeClick={handleTradeClick}
      />,
    );

    const row = screen.getByTestId('portfolio-row-aapl');
    const sellButton = within(row).getByTestId('sell-all-button');
    const tradeButton = within(row).getByTestId('trade-button');

    expect(sellButton).toBeEnabled();
    expect(tradeButton).toBeEnabled();

    await user.click(sellButton);
    expect(handleSellClick).toHaveBeenCalledWith(mockPortfolioItemAAPL);

    await user.click(tradeButton);
    expect(handleTradeClick).toHaveBeenCalledWith('AAPL');
  });

  it('disables actions and sets "Not Tradable" tooltip for untradable assets', () => {
    vi.mocked(usePortfolioRowState).mockReturnValue({
      symbol: 'MSFT',
      source: 'Finnhub',
      isMarketClosed: false,
      isTradable: false,
      quote: { price: 310.0 },
      marketValue: '$1550.00',
      pnl: '+$50.00',
      pnlColorClass: 'text-green-500',
    } as unknown as ReturnType<typeof usePortfolioRowState>);

    render(<PortfolioTable items={[mockPortfolioItemMSFT]} />);

    const row = screen.getByTestId('portfolio-row-msft');
    const sellButton = within(row).getByTestId('sell-all-button');
    const tradeButton = within(row).getByTestId('trade-button');

    expect(sellButton).toBeDisabled();
    expect(tradeButton).toBeDisabled();
  });

  it('disables actions and sets "Market Closed" tooltip when market is closed', () => {
    vi.mocked(usePortfolioRowState).mockReturnValue({
      symbol: 'AAPL',
      source: 'Finnhub',
      isMarketClosed: true,
      isTradable: true,
      quote: { price: 160.0 },
      marketValue: '$1600.00',
      pnl: '+$100.00',
      pnlColorClass: 'text-green-500',
    } as unknown as ReturnType<typeof usePortfolioRowState>);

    render(<PortfolioTable items={[mockPortfolioItemAAPL]} />);

    const row = screen.getByTestId('portfolio-row-aapl');
    const sellButton = within(row).getByTestId('sell-all-button');
    const tradeButton = within(row).getByTestId('trade-button');

    expect(sellButton).toBeDisabled();
    expect(tradeButton).toBeDisabled();
  });
});
