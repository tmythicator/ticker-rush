import { renderHook } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { useDashboardStats } from './useDashboardStats';
import { mockUserParticipating } from '@/test/mocks';

describe('useDashboardStats', () => {
  it('returns placeholder values when user is null', () => {
    const { result } = renderHook(() => useDashboardStats(null));

    expect(result.current).toHaveLength(3);
    expect(result.current[0]).toMatchObject({ label: 'Cash Balance', value: '--' });
    expect(result.current[1]).toMatchObject({ label: 'Invested Capital', value: '$0.00' });
    expect(result.current[2]).toMatchObject({ label: 'Open Positions', value: '0' });
  });

  it('calculates stats correctly based on user state', () => {
    const { result } = renderHook(() => useDashboardStats(mockUserParticipating));

    expect(result.current).toHaveLength(3);
    // mockUserParticipating balance = 5000
    expect(result.current[0]).toMatchObject({ label: 'Cash Balance', value: '$5000.00' });
    // mockUserParticipating portfolio has 10 AAPL @ 150 = 1500
    expect(result.current[1]).toMatchObject({ label: 'Invested Capital', value: '$1500.00' });
    expect(result.current[2]).toMatchObject({ label: 'Open Positions', value: '1' });
  });
});
