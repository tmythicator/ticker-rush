import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useQuotesSSE } from './useQuotesSSE';
import { useQuoteQuery } from '@/hooks/useQuoteQuery';
import { useQueryClient } from '@tanstack/react-query';
import { mockActiveQuote } from '@/test/mocks';

vi.mock('@/hooks/useQuoteQuery', () => ({
  useQuoteQuery: vi.fn(),
}));

vi.mock('@tanstack/react-query', () => ({
  useQueryClient: vi.fn(),
}));

class MockEventSource {
  url: string;
  onopen: ((event: Event) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  listeners: Record<string, ((event: MessageEvent) => void)[]> = {};
  close = vi.fn();

  constructor(url: string) {
    this.url = url;
    MockEventSource.instances.push(this);
  }

  addEventListener(type: string, listener: (event: MessageEvent) => void) {
    if (!this.listeners[type]) {
      this.listeners[type] = [];
    }
    this.listeners[type].push(listener);
  }

  removeEventListener(type: string, listener: (event: MessageEvent) => void) {
    if (this.listeners[type]) {
      this.listeners[type] = this.listeners[type].filter((l) => l !== listener);
    }
  }

  triggerOpen() {
    if (this.onopen) this.onopen(new Event('open'));
  }

  triggerError() {
    if (this.onerror) this.onerror(new Event('error'));
  }

  triggerEvent(type: string, data: unknown) {
    const event = new MessageEvent(type, { data: JSON.stringify(data) });
    if (this.listeners[type]) {
      this.listeners[type].forEach((l) => l(event));
    }
  }

  static instances: MockEventSource[] = [];
}

describe('useQuotesSSE', () => {
  let queryClientMock: { setQueryData: ReturnType<typeof vi.fn> };

  beforeEach(() => {
    vi.clearAllMocks();
    MockEventSource.instances = [];
    Object.defineProperty(globalThis, 'EventSource', {
      value: MockEventSource,
      configurable: true,
      writable: true,
    });

    queryClientMock = {
      setQueryData: vi.fn(),
    };
    (useQueryClient as unknown as ReturnType<typeof vi.fn>).mockReturnValue(queryClientMock);
    (useQuoteQuery as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockActiveQuote,
    });
  });

  afterEach(() => {
    Object.defineProperty(globalThis, 'EventSource', {
      value: undefined,
      configurable: true,
      writable: true,
    });
  });

  it('does not open EventSource if symbol is null', () => {
    renderHook(() => useQuotesSSE(null));
    expect(MockEventSource.instances).toHaveLength(0);
  });

  it('subscribes to SSE events and updates query cache', () => {
    const { result } = renderHook(() => useQuotesSSE('AAPL'));

    expect(MockEventSource.instances).toHaveLength(1);
    const instance = MockEventSource.instances[0];
    expect(instance.url).toContain('/v1/quotes/events?symbol=AAPL');

    expect(result.current.quote).toEqual(mockActiveQuote);
    expect(result.current.error).toBeNull();

    // Trigger SSE event
    const newQuote = { ...mockActiveQuote, price: 160, timestamp: '2026-06-14T20:00:00Z' };
    act(() => {
      instance.triggerEvent('quote', newQuote);
    });

    expect(queryClientMock.setQueryData).toHaveBeenCalled();
  });

  it('handles connection error state correctly', () => {
    const { result } = renderHook(() => useQuotesSSE('AAPL'));
    const instance = MockEventSource.instances[0];

    act(() => {
      instance.triggerError();
    });

    expect(result.current.error).toBeInstanceOf(Event);
  });
});
