export const queryConfig = {
  quotes: {
    // Keep quote cache fresh for 30 seconds
    staleTime: 1000 * 30,
    // Disable aggressive refetching on window focus/mount/reconnect for quotes
    // because real-time SSE updates are already pushed dynamically.
    refetchOnWindowFocus: false,
    refetchOnMount: false,
    refetchOnReconnect: false,
  },
  history: {
    // Keep historical chart data fresh for 5 minutes
    staleTime: 1000 * 60 * 5,
    refetchOnWindowFocus: false,
  },
  homeChart: {
    // Keep home chart history fresh for 3 minutes
    staleTime: 1000 * 60 * 3,
  },
  user: {
    // Keep user profiles fresh for 5 minutes
    staleTime: 1000 * 60 * 5,
  },
} as const;
